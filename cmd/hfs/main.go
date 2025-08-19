package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/c0va23/go-proxyprotocol"
	"github.com/rs/cors" // AJOUTÉ : Importer la bibliothèque CORS
	"github.com/zhangyoufu/hfs"
)

// La fonction umask est nécessaire pour compiler votre code original sous non-Linux.
func umask(mask int) int {
	return 0 // Implémentation factice pour Windows/autres OS
}

func main() {
	var (
		network     string
		address     string
		root        string
		enableProxy bool
		accessLog   bool
		dirList     bool
		dirFirst    bool
		ignoreCase  bool
		dotFile     bool
		indexPage   bool
	)

	flag.StringVar(&network, "network", "tcp", "tcp or unix")
	flag.StringVar(&address, "addr", "localhost:8000", "address:port for tcp, or filesystem path for unix")
	flag.StringVar(&root, "root", ".", "`path` of document root")
	flag.BoolVar(&enableProxy, "proxy", false, "enable PROXY protocol support")
	flag.BoolVar(&accessLog, "log", true, "enable access log")
	flag.BoolVar(&dirList, "list", true, "enable directory listing")
	flag.BoolVar(&dirFirst, "dirfirst", true, "list directories before files")
	flag.BoolVar(&ignoreCase, "igncase", true, "case insensitive sorting")
	flag.BoolVar(&dotFile, "dotfile", false, "enable listing and serving dot files (default false)")
	flag.BoolVar(&indexPage, "index", true, "enable serving index.html")
	flag.Parse()

	sorter := hfs.Sorter(nil)
	if dirList {
		sorter = hfs.NaiveSorter{
			DirectoryFirst: dirFirst,
			IgnoreCase:     ignoreCase,
		}
	}

	var logger *log.Logger
	if accessLog {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	switch network {
	case "tcp":
	case "unix":
		_ = os.Remove(address)
		_ = umask(0o007)
	default:
		log.Fatal("unsupported network")
	}

	rawListener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("unable to listen: ", err)
	}
	ln := rawListener.(net.Listener)

	if enableProxy {
		proxyListener := proxyprotocol.NewDefaultListener(ln)
		ln = proxyListener
	}

	// MODIFIÉ : Créer le handler hfs d'abord
	hfsHandler := &hfs.FileServer{
		FileSystem:       http.Dir(root),
		Sorter:           sorter,
		AccessLog:        logger,
		ErrorLog:         log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		DirectoryListing: dirList,
		ServeDotFile:     dotFile,
		ServeIndexPage:   indexPage,
	}

	// AJOUTÉ : Configurer le middleware CORS
	corsHandler := cors.New(cors.Options{
		// Pour la production, remplacez "*" par vos domaines autorisés.
		// Exemple : []string{"https://mon-site.com", "https://mon-autre-site.com"}
		AllowedOrigins: []string{"*"},

		// Méthodes HTTP autorisées. Pour un serveur de fichiers, GET et HEAD suffisent.
		// OPTIONS est géré automatiquement par la bibliothèque pour les requêtes "preflight".
		AllowedMethods: []string{"GET", "HEAD", "OPTIONS"},

		// Permet au navigateur d'envoyer des cookies ou des en-têtes d'authentification.
		AllowCredentials: true,

		// En-têtes autorisés
		AllowedHeaders: []string{"*"},
	}).Handler(hfsHandler) // <-- Encapsule le handler hfs original

	server := http.Server{
		// MODIFIÉ : Utiliser le handler encapsulé par CORS
		Handler: corsHandler,
		DisableGeneralOptionsHandler: true,
	}

	log.Printf("Serving %s on [%s] %s with CORS enabled", root, network, address)
	log.Fatal(server.Serve(ln))
}
