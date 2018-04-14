# pass2pass

Tool for migrate passwords from one password manager format to another.

I created this tool for my own needs, use it on your own risk.

There are some other tools for this task, such as:

https://github.com/roddhjav/pass-import#readme
https://git.zx2c4.com/password-store/tree/contrib/importers/lastpass2pass.rb

Try them, they are simpler, better tested and may suit your needs.

This one currently can only parse passwords, exported into CSV format by
LastPass and store them either in plain-text files (for debugging) or using
`pass`/`gopass` commands (see https://www.passwordstore.org/ and
https://www.justwatch.com/gopass/). 

I started fiddling with my own tool because I want to use different scheme,
than mentioned tools are producing. Namely, it's a directory-based schema,
which described as an "alternative approach" on the `pass`'s website.

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

Command to store password also can be supplied via CLI, so it should be
possible to use some other CLI password manager to import passwords, but I
don't need it currently, so I haven't tested it with anything except `gopass`.

I tried to make it extensible just in case I (or may be someone else) will
need to add support of other password managers into it. It should be enough to
add parser/store code into corresponded packages by analogy with existing one.
