/*
Package clippy handles registering CLI flags via a plug-in mechanism.
Additionally, optional CLI flag handling can be carried out after the CLI flags
have been parsed and just right before the selected command is about to run.

Packages providing CLI flags need to register themselves using the go-plugger
mechanism providing functions assignable to [cliplugin.SetupCLI]. This makes it
possible to register multiple flags from the same package in a highly modular
fashion just by specifying individual plug-in functions, for instance, per each
individual flag.

Registering [cliplugin.BeforeCommand] functions allows for modular processing of
flags before the (root) command is run.

Registration of the exported functions should be done in init functions. For
better modularity, multiple registration-related init functions can perfectly
co-exist within the same package. You might want to specify different “plug-in
names” upon function registration.

It already sufficies when a cmd package references a CLI-related package to pull
in its CLI flag registrations. A cmd package then should make sure to call
[AddFlags] and [BeforeCommand] respectively.

AddFlags should be called after your cmd package has created the root command
object and is ready for registering flags.

BeforeCommand should be called from the PersistentPreRunE hook function of your
cmd package's root command object.
*/
package clippy
