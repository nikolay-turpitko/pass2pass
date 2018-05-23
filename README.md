# pass2pass

## Intro

Tool for migrate passwords from one password manager format to another.

I created this tool for my own needs, use it on your own risk.

There are some other tools for this task, such as:

- https://github.com/roddhjav/pass-import#readme
- https://git.zx2c4.com/password-store/tree/contrib/importers/lastpass2pass.rb

Try them, they are simpler, better tested and may suit your needs.

This one currently can only parse passwords, exported into CSV format by
LastPass and store them either in plain-text files (for debugging) or using
`pass`/`gopass` commands (see https://www.passwordstore.org/ and
https://www.justwatch.com/gopass/). 

I started fiddling with my own tool because I want to use different scheme,
than mentioned tools are producing. Namely, it's a directory-based scheme,
which described as an "alternative approach" on the `pass`'s website.

_Update_: later I decided to switch to scheme, supported by `gopass`'s companion
browser plugin. Changing layout schema is much simple with template-based tool.

To have some more fun during this tool development and to sharpen my own skills
I added some concurrency to it (not like it's very necessary here). Actually, 
concurrency even cause errors with `gopass` when `gopass` tries to perform `git 
commit`, so I had to wrap `gopass` invocation with `flock`.

Because of my particular requirements and for more flexibility, tool supports
pre-processing password entry paths with custom external commands/scripts
(bash/sed/awk/etc). It made the tool a bit slower than its initial native
implementation, but this is completely optional, fast enough and quite
convenient. There is a room for improvement, if anyone ever need it:  currently
processes for commands are created and deleted in the loop for every password
entry. It's possible to refactor this code so that they were created only once
at program startup and data was piped through them via their stdin/stdout. But
I'm not currently in the mood to make this optimization.

_Update_: I added an ability to use Go templates instead of external commands.
This a bit more limited (to only those Go template functions I bothered to
implement for myself), but this is a bit faster then invoking external scripts
(which is still possible). So, currently performance is restricted only by
file operations (reading of templates) and invoking `gopass` command.

Command to store password also can be supplied via CLI, so it should be
possible to use some other CLI password manager to import passwords, but I
don't need it currently, so I haven't tested it with anything except `gopass`.

I tried to make it extensible just in case I (or may be someone else) will
need to add support of other password managers into it. It should be enough to
add parser/store code into corresponded packages by analogy with existing one.

## Usage

I created it for my own use, so it's not quite user friendly. It's not `go
get`able.  Use git to clone the project, then `cp ./env-nick ./env-$USER` and
modify it, then use `./run` to have an idea how to invoke it. And yes, you'll
need Go compiler installed. And don't forget to execute `dep ensure` in 
`./src/pass2pass` dir. After that you should be able to build and run it. To 
see help on flags and usage you may want to run it with `-h` flag. And, of 
course, you may want to look into sources.  OK, if someone else ever need it, 
I'll think about making it more convenient to install.

Take a look into `./scriptlets` and `./templates` folders. These are for user's
customization. It's where I put my own scripts and templates I used to
transform pass entries during migration.

## Troubleshooting

In case of issues with GPG pinentry (login) dialog try to use `pass show
<password>` to show some existing password, so that pass invoke GPG and ask for
it's password.  It should cache password after that and shouldn't ask it during
`pass2pass` invocation.

## Status

Currently it's completely for my own needs. Use it for your own risk. You may
fork it and brush up if you like. Ask questions, report bugs in Github's issue
tracker, send PRs. If I'll see some interest to it I'll, probably, find some
time to make it more user friendly. For my own needs it's perfect as it is.

## Links

- https://cerb.ai/guides/mail/gpg-setup-on-mac/#installing-software
- https://github.com/justwatchcom/gopass/blob/master/docs/setup.md#filling-in-passwords-from-browser
- https://github.com/martinhoefling/gopassbridge

