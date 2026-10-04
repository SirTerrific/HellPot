# Fichiers

- [CI, sécurité et releases](ci-securite-et-releases.md) - Les trois workflows GitHub Actions du dépôt, les contrôles de sécurité automatiques, la mise à jour des dépendances par Dependabot, la procédure de publication par tag et la règle de nommage des versions.
- [Image Docker et publication sur GHCR](docker-et-publication.md) - Comment l'image HellPot est construite (Dockerfile multi-étapes, compilation croisée amd64 et arm64), ce qu'elle contient, comment elle est publiée sur GHCR et comment la déployer avec docker run ou docker compose.
- [Journalisation et format des logs](journalisation.md) - Comment HellPot écrit ses logs (JSON zerolog dans un fichier ou sur stdout, console lisible), le nom du fichier, le catalogue des messages et de leurs champs, et la contrainte de ne pas changer ce format.
- [Limites de connexions et performance](limites-et-performance.md) - Les réglages qui bornent la charge de HellPot (connexions par IP, concurrence, timeouts), l'effet d'un reverse proxy sur la limite par IP, ce qui a été mesuré, et pourquoi le pool de buffers a été retiré.
