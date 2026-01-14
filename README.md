# Pleesah havnesjef

Denne har ansvaret for å administrere teams.

Den setter opp følgende for hvert team:

- Namespace
- Service account
- Secret
- Rolebinding

Når dette er gjort lager den en token for service accounten og viser en `KUBECONFIG` som teamet kan bruke.
