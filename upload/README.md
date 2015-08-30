upload
======

**upload** - files upload service in golang.

[![last-version-blue](https://cloud.githubusercontent.com/assets/6298396/5602522/8967405e-935b-11e4-8777-de3623ed6ad7.png)](https://github.com/geosoft1/upload/archive/master.zip)

![screenshot from 2014-08-31 12 02 23](https://cloud.githubusercontent.com/assets/6298396/4101485/cfd8bb74-30ed-11e4-8061-75ac3df336a1.png)

Useful if you want to send big files to others and common mail servers don't accept
files over a standard length (usually a few MB).

* **Using**
    * Compile. See [here](https://github.com/geosoft1/tools) how.
    * Copy `` upload `` folder in your account on your Linux machine.
    * Run  `` ./upload `` file from upload folder.
    * Open [http://localhost:8080](http://localhost:8080) in browser. Replace localhost with server ip if you access upload service from another computer. If you are behind a router remember to forward 8080 port. Change this port as needed.
    * Use user and password from `` secrets `` file. Change this file as needed.	
    * Upload a file and send received link to others for download.

Files are transfered in [/files](http://localhost:8080/files) folder.

Automaticaly start service by adding in /etc/rc.local

     sudo nano /etc/rc.local
	
and add following line before `` exit ``

     /home/user/cws/upload/upload &

Use cron to keep server clean

     crontab -e

and add following line to clean files older than 24h

    MAILTO=""
    * * * * * find /home/user/cws/upload/files/* -mtime +1 -exec rm {} \;

Change **user** with your linux user. **cws** means custom web services and it's my convention. clear MAILTO to avoid sending mails to root if no files to delete.
Feel free to use any location.

No root rights are required.
