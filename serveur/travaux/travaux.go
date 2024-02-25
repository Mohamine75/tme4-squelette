package travaux

import (
	"math/rand"

	st "tme4-squelette/client/structures"
)

// *** LISTES DE FONCTION DE TRAVAIL DE Personne DANS Personne DU SERVEUR ***
// Essayer de trouver des fonctions *différentes* de celles du client

func f1(p st.Personne) st.Personne {
	// A FAIRE
	np := p
	if np.Sexe == "M" {
		np.Prenom = "M." + p.Prenom
	} else {
		np.Prenom = "Mme." + p.Prenom
	}
	return np
}

func f2(p st.Personne) st.Personne {
	// A FAIRE
	np := st.Personne{
		Nom:    p.Nom,
		Prenom: p.Prenom,
		Age:    18,
		Sexe:   p.Sexe,
	}
	return np
}

func f3(p st.Personne) st.Personne {
	// A FAIRE
	np := st.Personne{
		Nom:    "KESSAL",
		Prenom: p.Prenom,
		Age:    22,
		Sexe:   "M",
	}
	return np
}

func f4(p st.Personne) st.Personne {
	// A FAIRE
	np := st.Personne{
		Nom:    "ZEMALI",
		Prenom: p.Prenom,
		Age:    22,
		Sexe:   "F",
	}
	return np
}

func UnTravail() func(st.Personne) st.Personne {
	tableau := make([]func(st.Personne) st.Personne, 0)
	tableau = append(tableau, func(p st.Personne) st.Personne { return f1(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f2(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f3(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f4(p) })
	i := rand.Intn(len(tableau))
	return tableau[i]
}
