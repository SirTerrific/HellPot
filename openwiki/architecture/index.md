# Fichiers

- [Moteur de chaînes de Markov (heffalump)](moteur-markov-heffalump.md) - Comment le package heffalump fabrique le flux de texte infini envoyé aux robots, de la décompression du texte source à l'écriture dans la réponse HTTP, et comment il s'arrête quand le client part.
- [Serveur HTTP, routage et piège](serveur-http-et-routage.md) - Comment HellPot configure son serveur fasthttp, décide quelles requêtes tombent dans le piège (catchall, paths, robots.txt, liste noire de user agents), identifie le client et écoute en TCP ou sur socket Unix.
- [Vue d'ensemble de l'architecture](vue-densemble.md) - Les cinq packages de HellPot, leur rôle, et la séquence de démarrage du processus, de l'initialisation de la configuration et du logger jusqu'au serveur HTTP et à l'arrêt sur signal.
