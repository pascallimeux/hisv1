get consent from github
wget https://raw.githubusercontent.com/pascallimeux/consentSC/master/consentv2/consentv2.go 
wget https://raw.githubusercontent.com/pascallimeux/consentSC/master/consentv2/consentv2_test.go
go test
chmod 755 ../consentv2/
chmod 644 *.* 

Dans les constantes on va définir les index et les libellés des erreurs
les index permettent les recherches multi critères dans la BDD

on va créer les index

on définit consentCC (la classe)
on définit la structure consent qui va représenter un consentement

on doit implementer init et invoke (init utiliser à l'initialisation du CC, et invoke sera le point d'entrée du CC)
c'est dans la méthode invoke que nous allons diriger les requetes vers les fonctions du CC.

ensuite on implemente chaque fonction dans lesquelles on récupère les arguments et on retourne une erreur ou un shimsuccess
avec l'objet sous forme de bytes.

getConsentsByIndex pour effectuer des recherches multi critères dans le BDD