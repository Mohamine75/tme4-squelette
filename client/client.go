package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
	st "tme4-squelette/client/structures"
	"tme4-squelette/client/travaux"
)

var ADRESSE string = "localhost"                       // adresse de base pour la Partie 2
var FICHIER_SOURCE string = "./elus-municipaux-cm.csv" // fichier dans lequel piocher des personnes
var TAILLE_SOURCE int = 450000                         // inferieure au nombre de lignes du fichier, pour prendre une ligne au hasard
const TAILLE_G int = 5                                 // taille du tampon des gestionnaires
const NB_G int = 2                                     // nombre de gestionnaires
var NB_P int = 2                                       // nombre de producteurs
var NB_O int = 2                                       // nombre d'ouvriers
var NB_PD int = 2                                      // nombre de producteurs distants pour la Partie 2
var cpt_id = 0

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"} // une personne vide

type tuple struct {
	ligne      int
	retourChan chan string
}

type tuple_proxy struct {
	methode string
	id      int
	retour  chan string
}

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int
type personne_emp struct {
	// A FAIRE
	numero_ligne int
	statut       string
	personne     st.Personne
	afaire       []func(personne st.Personne) st.Personne
	ligne_retour chan tuple
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	//A Faire
	identifiant int
	proxy_can   chan tuple_proxy
}

/*func Newpersonne_emp() personne_emp {
	paquet := personne_emp{}
	paquet.statut = "V"
	paquet.afaire = []func(){}
	rand.Seed(time.Now().Unix())
	paquet.numero_ligne = rand.Intn(TAILLE_SOURCE)
	return paquet
}*/

// interface des personnes manipulees par les ouvriers, les
type personne_int interface {
	initialise()          // appelle sur une personne vide de statut V, remplit les champs de la personne et passe son statut à R
	travaille()           // appelle sur une personne de statut R, travaille une fois sur la personne et passe son statut à C s'il n'y a plus de travail a faire
	vers_string() string  // convertit la personne en string
	donne_statut() string // renvoie V, R ou C
}

// fabrique une personne à partir d'une ligne du fichier des conseillers municipaux
// à changer si un autre fichier est utilisé
func personne_de_ligne(l string) st.Personne {
	separateur := regexp.MustCompile(",") // oui, les donnees sont separees par des tabulations ... merci la Republique Francaise
	separation := separateur.Split(l, -1)
	naiss, _ := time.Parse("2/1/2006", separation[9])
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	return st.Personne{Nom: separation[6], Prenom: separation[7], Sexe: separation[8], Age: agec}
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

func (p *personne_emp) initialise() {
	ch_retour := make(chan string)
	tmp := tuple{ligne: p.numero_ligne, retourChan: ch_retour}
	p.ligne_retour <- tmp
	ligne_string := <-ch_retour

	p.personne = personne_de_ligne(ligne_string)
	rand.Seed(time.Now().Unix())
	nb_alea_funs := rand.Intn(5) + 1
	for i := 0; i < nb_alea_funs; i++ {
		p.afaire = append(p.afaire, travaux.UnTravail())
	}
	p.statut = "R"
}

func (p *personne_emp) travaille() {
	/**
	On ne devrait se trouver dans aucun de ces 3 cas, mais par sécurité et debuggage on le met
	comme des pré-cond
	*/
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

func (p *personne_emp) vers_string() string {
	res := "Nom : " + p.personne.Nom + "\n Prenom : " + p.personne.Prenom + " \n Age : " + fmt.Sprint(p.personne.Age) + "\n Sexe : " + p.personne.Sexe
	return res
}

func (p *personne_emp) donne_statut() string {
	// A FAIRE
	return p.statut
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES DISTANTES (PARTIE 2) ***
// ces méthodes doivent appeler le proxy (aucun calcul direct)

func (p personne_dist) initialise() {
	// A FAIRE
	message_proxy := tuple_proxy{
		methode: "initialise",
		id:      p.identifiant,
	}
	p.proxy_can <- message_proxy
}

func (p personne_dist) travaille() {
	// A FAIRE
	message_proxy := tuple_proxy{
		methode: "travaille",
		id:      p.identifiant,
	}
	p.proxy_can <- message_proxy
}

func (p personne_dist) vers_string() string {
	// A FAIRE
	message_proxy := tuple_proxy{
		methode: "vers_string",
		id:      p.identifiant,
		retour:  make(chan string),
	}
	p.proxy_can <- message_proxy
	return <-message_proxy.retour
}

func (p personne_dist) donne_statut() string {
	// A FAIRE

	message_proxy := tuple_proxy{
		methode: "donne_statut",
		id:      p.identifiant,
		retour:  make(chan string),
	}
	p.proxy_can <- message_proxy
	return <-message_proxy.retour
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy(port string, proxy_can chan tuple_proxy) {
	// A FAIRE
	/**
	 */
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal("Erreur lors de la connexion au serveur:", err)
	}
	defer conn.Close()
	for {
		tuple := <-proxy_can
		methode := tuple.methode
		id := tuple.id

		// Construire le message à envoyer au serveur
		message := fmt.Sprintf("%d,%s\n", id, methode)
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Erreur lors de l'envoi du message au serveur:", err)
			continue
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Erreur lors de l'envoi du message au serveur:", err)
			continue
		}

		rep := string(buffer[:n])
		fmt.Println("Reponse du serveur : ", rep)
		tuple.retour <- rep
	}
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur(data chan tuple) {
	for {
		fichier, err := os.Open(FICHIER_SOURCE)
		if err != nil {
			log.Fatal(err)
		}
		defer fichier.Close()
		scanner := bufio.NewScanner(fichier)
		_ = scanner.Scan()
		msg := <-data
		var i = 0
		/*
			Pour accéder à la i-ème ligne nous avions penser à parser le texte par \n, mais ça ne marchait pas
		*/
		for scanner.Scan() {
			ligne := scanner.Text()
			if i == msg.ligne {
				msg.retourChan <- ligne
				break
			}
			i++
		}
	}
}

// Partie 1: récupèrent des personne_int depuis les gestionnaires, font une opération dépendant de donne_statut()
// Si le statut est V, ils initialise le paquet de personne puis le repasse aux gestionnaires
// Si le statut est R, ils travaille une fois sur le paquet puis le repasse aux gestionnaires
// Si le statut est C, ils passent le paquet au collecteur
func ouvrier(gestionnaires chan personne_int, collecteur chan personne_int) {
	// A FAIRE
	for {
		pers := <-gestionnaires
		if pers.donne_statut() == "V" {
			pers.initialise()
			gestionnaires <- pers
		} else if pers.donne_statut() == "R" {
			pers.travaille()
			gestionnaires <- pers
		} else if pers.donne_statut() == "C" {
			collecteur <- pers
		}
	}
}

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(ligne_retourChan chan tuple, gestionnaire_chan chan personne_int) {
	// A FAIRE
	for {
		random := rand.Intn(TAILLE_SOURCE)
		p := personne_emp{
			numero_ligne: random,
			statut:       "V",
			personne:     pers_vide,
			afaire:       []func(personne st.Personne) st.Personne{},
			ligne_retour: ligne_retourChan,
		}

		gestionnaire_chan <- personne_int(&p)

	}
}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V
func producteur_distant(gestionnaire chan personne_int, proxy_channel chan tuple_proxy) {
	// A FAIRE
	for {
		meth := "creer"
		p := personne_dist{identifiant: cpt_id, proxy_can: proxy_channel}
		proxy_channel <- tuple_proxy{id: cpt_id, methode: meth}
		gestionnaire <- p
		cpt_id++
	}
}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets surlesquels ils travaillent pour en prendre des autres
func gestionnaire(producteurs chan personne_int, ouvriers chan personne_int) {
	// A FAIRE
	file := make([]personne_int, 0)
	// On laisse toujours un emplacement libre pour les ouvriers, pour que si il ne reste qu'un emplacement, on
	// libère d'abord un ouvrier, sinon famine pour ouvrier car les producteurs ont rempli la file
	for {
		if len(file) == TAILLE_G-1 {
			pers_tmp := <-ouvriers
			file = append(file, pers_tmp)
		} else if len(file) == TAILLE_G {
			// Ne rien faire car file pleine
			continue
		} else {
			select {
			case pers_tmp := <-producteurs:
				file = append(file, pers_tmp)
			case pers_tmp := <-ouvriers:
				file = append(file, pers_tmp)
			}
		}
		go func() {
			ouvriers <- file[0]
			file = file[:1]

		}()
	}
}

// Partie 1: le collecteur recoit des personne_int dont le statut est c, il les collecte dans un journal
// quand il recoit un signal de fin du temps, il imprime son journal.
func collecteur(collecteur chan personne_int, fintemps chan int) {
	// A FAIRE
	str_res := ""
	for {
		select {
		case personne := <-collecteur:
			str_res += personne.vers_string() + "\n"
		case _ = <-fintemps:
			fmt.Println(str_res)
			fintemps <- 0
		}

	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire
	if len(os.Args) < 2 {
		fmt.Println(" Format: client <port> <millisecondes d'attente>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1]) // utile pour la partie 2
	tempsDeSleep := os.Args[2]
	duree, _ := time.ParseDuration(tempsDeSleep + "s")
	fintemps := make(chan int)
	// A FAIRE
	// creer les canaux
	can_gest := make(chan personne_int)
	can_ouvrier := make(chan personne_int)
	can_lecteur := make(chan tuple)
	can_collecteur := make(chan personne_int)
	can_proxy := make(chan tuple_proxy)
	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	go func() {
		lecteur(can_lecteur)
	}()
	go func() {
		collecteur(can_collecteur, fintemps)
	}()
	for i := 0; i < NB_G; i++ {
		go func() {
			gestionnaire(can_gest, can_ouvrier)
		}()
	}
	for i := 0; i < NB_P; i++ {
		go func() {
			producteur(can_lecteur, can_gest)
		}()

	}
	for i := 0; i < NB_O; i++ {
		go func() {
			ouvrier(can_ouvrier, can_collecteur)
		}()
	}

	// lancer les goroutines (partie 2): des producteurs distants, un proxy
	go func() {
		proxy(strconv.Itoa(port), can_proxy)
	}()
	for i := 0; i < NB_PD; i++ {
		go func() {
			producteur_distant(can_gest, can_proxy)
		}()
	}
	time.Sleep(duree)
	fintemps <- 0
	<-fintemps
	return
}
