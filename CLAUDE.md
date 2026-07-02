# CLAUDE.md — clearline-crm

<!-- rooche-team-conventions -->
## Rooche team conventions

Shared rules for this repo (one of ~30 in the Rooche / `wooshdigital` set,
maintained by a very small team). Any repo-specific notes go elsewhere in this file.

### Deploy
When the user says **"deploy"**:
1. **Merge** the current work into the **`main`** branch and push to origin.
2. **Pull on the respective server** — SSH to the box that hosts this repo and
   `git pull --ff-only` in its deploy directory (then restart the service if it
   needs it: `pm2 restart …` / `docker compose up -d …`).

The deploy target (host + path) for each repo is recorded in the rd-tools
`scripts/repos-lib.mjs` `DEPLOYMENT` map. Hosts: main box `159.65.11.83`,
auth/IdP droplet `143.198.95.106`.

### Config & secrets — committed to the repo (deliberate)
Small team + 30+ repos → **env/config files are committed to git** (not
gitignored) so any maintainer can clone-and-run. Naming:
- **`env_live`** — production values
- **`env_staging`** — staging values (once a staging environment exists)
- **`env_local_<username>`** — a maintainer's local overrides (e.g. `env_local_rj`)

**Mobile apps:** commit the **signing key / keystore and its password(s)** into
the repo as well (keystore file + the store/key passwords in `env_live`).

> Accepted tradeoff: if this repo is public, committed secrets are world-readable
> and remain in git history permanently. Rotate anything that leaks.
