# <b>C</b>onfiguration <b>D</b>ata<b>B</b>ase (CDB)

The user friendly configuration management database that integrates with your
DevOps tooling.

**Note:** CDB is not currently suitable for storing secrets, it is recommended
to use the secrets manager provided by your cloud provider or a self-hosted
alternative. We may support this in the future.

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
