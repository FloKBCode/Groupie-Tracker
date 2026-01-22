# **ARCHITECTURE DU PROGRAMME**

### API

L'API se compose uniquement du fichier client.go, qui utilise des fonctions effectuant des requetes HTTP Get/vérifier le statut
des requetes, et gérer des erreurs de timeout.

### MODELS

Ce dossier est réservé au stockage de données, composé de structures. Toutes ces données sont utilisées par les fonctions utilitaires
dans le dossier ui et services, notamment search.go ou fetch.go par exemple.

### SERVICES

Services est un ensemble qui rassemble toutes les fonctions qui servent au fonctionnement du back-end de l'application. Les données
contenues dans models sont exécutées via les fonctions de service.


fetch.go récupère des valeurs tel que des noms d'artistes ou des dates de représentations via l'api.

geocoding.go et gedocoding_preloader.go sont des fichiers avec des fonctions utilitaires afin de gérer la geolocalisation de
l'utilisateur. 
geocoding.go fonctionne en convertissant des localisations textuelles en coordoonées GPS grace a l'API.

fuzzy_search.go est un fichier qui permet de gérer la "recherche floue", c'est a dire le programme qui corrige ce que l'utilisateur a
écrit en lui proposant un résultat similaire a ce qu'il a écrit. (Cela est géré par la "distance de Levenshtein.")

### UI

UI compose l'entièreté des éléments qui gérènt le front-end du programme. Les éléments visuels, basés sur Fyne servent en tant
qu'alternative a du html plus classique.

app.go est le fichier le plus important du projet, puisque c'est lui qui permet d'afficher tout ce qu'on voit a l'écran, ainsi 
l'entièreté des fonctions du programme repose sur les ordres de app.go; Le programme exécute d'autres fonctions a son démarrage, tel que
l'affichage de la couleur de l'écran, le nettoyage de l'écran lors du changement de page, ou encore bien evidemment la gestion du 
démarrage de l'application.