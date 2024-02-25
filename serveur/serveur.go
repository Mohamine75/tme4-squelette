package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	st "tme4-squelette/client/structures"
	"tme4-squelette/serveur/travaux"
)

var ADRESSE = "localhost"

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "F"}
var table_association = make([]personne_serv, 0) //tableau d'association d'un identifiant de type entier(ici c'est l'indice) avec une personne_serv

// type d'un paquet de personne stocke sur le serveur, n'implemente pas forcement personne_int (qui n'existe pas ici)
type personne_serv struct {
	// A FAIRE
	identifiant int
	statut      string
	personne    st.Personne
	afaire      []func(personne st.Personne) st.Personne
}

// cree une nouvelle personne_serv, est appelé depuis le client, par le proxy, au moment ou un producteur distant
// produit une personne_dist
func creer(id int) *personne_serv {
	p := pers_vide
	fmt.Println("coucou")

	pers_serv := personne_serv{
		identifiant: id,
		statut:      "V",
		personne:    p,
		afaire:      []func(personne st.Personne) st.Personne{},
	}
	for len(table_association) <= id {
		table_association = append(table_association, pers_serv)
	}
	table_association[id] = pers_serv
	fmt.Println(table_association)
	return &pers_serv
}

// Méthodes sur les personne_serv, on peut recopier des méthodes des personne_emp du client
// l'initialisation peut être fait de maniere plus simple que sur le client
// (par exemple en initialisant toujours à la meme personne plutôt qu'en lisant un fichier)
func (p *personne_serv) initialise() {
	// A FAIRE
	p.personne = pers_vide
	rand.Seed(time.Now().Unix())
	nb_alea_funs := rand.Intn(5) + 1
	for i := 0; i < nb_alea_funs; i++ {
		p.afaire = append(p.afaire, travaux.UnTravail())
	}
	p.statut = "R"
}

func (p *personne_serv) travaille() {
	// A FAIRE
	if p.statut == "C" || p.statut == "V" || len(p.afaire) == 0 {
		panic("Probleme, aucun travail ne devrait être effectué")
	}
	p.personne = p.afaire[0](p.personne)
	if len(p.afaire) > 0 {
		p.afaire = p.afaire[1:]
	}
	if len(p.afaire) == 0 {
		p.statut = "C"
	}
}

func (p *personne_serv) vers_string() string {
	// A FAIRE
	res := "Nom : " + p.personne.Nom + "\n Prenom : " + p.personne.Prenom + " \n Age : " + fmt.Sprint(p.personne.Age) + "\n Sexe : " + p.personne.Sexe
	return res
}

func (p *personne_serv) donne_statut() string {
	// A FAIRE
	return p.statut
}

// Goroutine qui maintient une table d'association entre identifiant et personne_serv
// il est contacté par les goroutine de gestion avec un nom de methode et un identifiant
// et il appelle la méthode correspondante de la personne_serv correspondante
func mainteneur(id int, methode string) string {
	// A FAIRE
	switch methode {
	case "creer":
		creer(id)
		return "ok"
	case "initialise":
		pers_serv := table_association[id]
		pers_serv.initialise()
		return "ok"
	case "travaille":
		pers_serv := table_association[id]
		pers_serv.travaille()
		return "ok"
	case "vers_string":
		pers_serv := table_association[id]
		return pers_serv.vers_string()
	case "donne_statut":
		pers_serv := table_association[id]
		return pers_serv.donne_statut()
	default:
		creer(id)
		return "ok" // case default
	}
}

// Goroutine de gestion des connections
// elle attend sur la socketi un message content un nom de methode et un identifiant et appelle le mainteneur avec ces arguments
// elle recupere le resultat du mainteneur et l'envoie sur la socket, puis ferme la socket
func gere_connection(conn net.Conn) {
	// A FAIRE

	/**
	On recupere les message sur la socket, le mainteneur est simplement une fonction car on attend son return
	*/
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Erreur de lecture:", err)
			return
		}
		message := string(buffer[:n])

		// Extraction de l'identifiant et de la méthode à partir du message
		parts := strings.Split(message, ",")
		if len(parts) < 2 {
			fmt.Println("Message invalide:", message)
			return
		}
		id, err := strconv.Atoi(parts[0]) //on recupere l'identifiant
		if err != nil {
			fmt.Println("Identifiant invalide:", parts[0])
			return
		}

		methode := parts[1] //et on recupere la methode a donner au mainteneur

		// Vérification si l'identifiant existe dans la table d'association
		if len(table_association) < id {
			fmt.Println("Identifiant non trouvé:", id)
			return
		}
		res := mainteneur(id, methode)
		_, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("Erreur:", err)
			return
		}
		conn.Write([]byte(res))
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Format: client <port>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1]) // doit être le meme port que le client
	addr := ADRESSE + ":" + fmt.Sprint(port)
	// A FAIRE: creer les canaux necessaires, lancer un mainteneur
	ln, _ := net.Listen("tcp", addr) // ecoute sur l'internet electronique
	fmt.Println("Ecoute sur", addr)
	for {
		conn, _ := ln.Accept() // recoit une connection, cree une socket
		fmt.Println("Accepte une connection.")
		go func() { gere_connection(conn) }() // passe la connection a une routine de gestion des connections
		fmt.Println(conn)
	}

}

/**
La partie 2 nous a pris beaucoup de temps, nous pensons avoir bien fait les fonctions et tâches qui étaient demandées
dans le cahier des charges, le problème vient du parsing de méthode dans mainteneur qui semble défaillant à notre
grande surprise, la connection se fait.
*/
