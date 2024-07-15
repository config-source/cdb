# <b>C</b>onfiguration <b>D</b>ata<b>B</b>ase (CDB)

The user friendly configuration management database that integrates with your
DevOps tooling.

**Note:** CDB is not currently suitable for storing secrets, it is recommended
to use the secrets manager provided by your cloud provider or a self-hosted
alternative. We may support this in the future.

## The Problem

Everyone is lucky enough to have a test environment but not everyone is lucky
enough to have a separate production environment. For those of us lucky enough
to have separate environments configuration drift becomes a real issue.

You can solve this with a whole pile of YAML or `.tfvars` files, but that then
becomes its own problem of how to best manage them. This gets even worse if you
want to do something like [per pull-request Feature Environments](https://devops.com/implementing-an-on-demand-feature-environment/)
because it adds an [SRE](https://sre.google)-required element to manage.

If you do something fancier with your infrastructure as code to avoid the pile
of files then it becomes harder for non-SREs to make changes or customise their
environment. No matter what you do a Product Owner or QA person will almost
never feel comfortable interacting with it anyway.

Enter CDB an easy way to manage exactly the right amount of configuration and
expose it to anyone you want in your organisation (even those non-technical
folk).

## How It Works

CDB maintains primarily three objects that are all related:

1. Environments - These are logical containers for configuration and should map
   to your deployed environments and in CDB they will know their promotion
   relationship with each other.
2. Config Keys - These are the available keys to configure on environments,
   these are managed by a SysAdmin typically but you can change this with
   a setting so they are dynamically created.
3. Config Values - Instances of values for a config key and environment.

Since environments know who they promote to they will inherit configuration from
their parents if that configuration key has not been set on them directly. The
order of precedence is such that the nearest parent wins.

Consider a simple Dev -> Staging -> Production example:

![Environment Inheritance Diagram](/docs/images/environment-inheritance-diagram.png)

In this example the Production environment has three config values set on it:
`owner=SRE`, `minReplicas=10`, `maxReplicas=100`. Since Production is the
highest environment it inherits no values.

The Staging environment promotes to Production and so is a child environment of
Production. Since it does not have a direct value set for `owner` and
`maxReplicas` it inherits those from Production. However it does have
`minReplicas=1` set explicitly on it so that overrides the value from
Production.

The Dev environment promotes to Staging and has no direct Config Values set on
it so it inherits all of it's configuration. You can see that it inherits
`owner` and `maxReplicas` from Production as they aren't set on Staging but it
inherits `minReplicas` from Staging as Staging values take precedence over
Production as it is the nearer parent to Dev.

Whether a Config Key can be inherited or not can be configured on the Config
Keys themselves by an administrator. So for example if you don't want
`maxReplicas` to ever be inherited you can configure that on the Config Key
itself.

## Nifty Features

Beyond the configuration management database features CDB integrates with some devops tooling to solve some common pain points and add that little touch of convenience:

- (WIP) [A Terraform provider](https://github.com/config-source/terraform-provider-cdb) for querying
  configuration out of CDB.
- (TODO) It can be configured to manage ConfigMaps in Kubernetes and optionally
  restart specific deployments when they update
- (TODO) Can be configured to fire Webhooks at any URL when configuration
  updates with configurable JSON templates using the [Go templating](https://pkg.go.dev/text/template) 
  package.

## Contributing

You'll need [Docker](https://www.docker.com/get-started/) to effectively develop
on CDB as well as a text editor.

All development is done using Docker compose the best way to start a local
development server is to run `./scripts/dev-server` this will start everything
you need with docker compose and tail the logs of the server.

All other development tasks such as linting or tests have corresponding scripts
in `./scripts` which will run the appropriate commands inside of the docker
container so you can run them from your host machine.

To setup some starter data on your local instance simply run `./scripts/seed`
and that will wipe and load your database with:

- 4 environments (production, staging, dev1, dev2)
- 3 config keys (owner, minReplicas, maxReplicas)
- 5 config values
    - `owner=SRE`, `minReplicas=10`, and `maxReplicas=100` on production
    - `minReplicas=1` on production
    - `maxReplicas=10` on dev1
