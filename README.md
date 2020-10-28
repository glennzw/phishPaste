# Phish Paste
Phish Paste is a tiny utility to copy templates, pages, and sending profiles between different [Gophish](https://github.com/gophish/gophish) user accounts.

Example usage:
```
$ ./phishpaste --list-users
[+] Connecting to gophish.db (set $DATABASE_URL to use MySQL)
[+] Available users: 
	admin
	bobby
```
 
Let's copy all items from `admin` to `bobby`:

```
$ ./phishpaste --source admin --destination bobby --sending-profiles --email-templates --landing-pages 
[+] Connecting to gophish.db (set $DATABASE_URL to use MySQL)
[-] Dry run mode enabled, no data will actually be copied
[+] Copying 8 Landing Pages from 'admin' to 'bobby'
	✅ SimpleCreds 
	✅ HMRC 
	✅ Update Email Password 
	✅ Facebook 
	✅ Employee Evaluation 
	❌ Dropbox (destination user already has page of the same name) 
	✅ Antivirus Update 
	✅ Copy of Update Email Password 
[+] Copying 9 Email Templates from 'admin' to 'bobby'...
	✅ Template01 
	✅ HMRC 
	✅ Update Email Password 
	✅ MLTest01 
	✅ Facebook 
	✅ Employee Evaluation 
	❌ Dropbox (destination user already has template of the same name) 
	✅ Antivirus Update 
	✅ asriznet-test 
[+] Copying 9 Sending Profiles from 'admin' to 'bobby'...
	✅ SendGrid 
	✅ Mailslurp 
	✅ HMRC 
	✅ Update Email Password 
	✅ Facebook 
	✅ Employee Evaluation 
	❌ Dropbox (destination user already has sending profile of the same name) 
	✅ Antivirus Update 
	✅ Mailhog 

```

Note: If a destination item exists with the same name it's skipped (as you can see in the above output). The `--overwrite` option wil replace entrie with the same name.
