# Web Proxy Server

## Description
The web proxy is written in Go using built in and open source libraries. The program upon running will take input from the user through the CLI (Command Line Interface). The user may add or remove websites from the blacklist and view the current contents of the blacklist. The user will also see each time a request is allowed and if a request if blocked. A thread is running to periodically clear the cache of sites that have been on the cache for some time. The server is set up to listen to requests and the handle these requests using the inbuilt http handle function.

## The Cache
The server will cache any HTTP site that it requests and visits. The cache is stored as a map of the URL string as its keys and a custom website struct which contains the site responses and headers to send to the client if the site is requested again within the time limit allowed. If the user types /c They will be shown the cache information which shows the average time saved from using the cache and the bandwidth saved.

## The Blacklist
The server will store blacklisted domains in a map. To block a domain the user must enter /add <domain> in order to add a domain to the blacklist and they must enter /rmv <domain> to remove a domain from the blacklist. The user can also view the blacklist by entering /view. The proxy also blocks subdomains on the URL entered i.e. if the user blocked tcd.ie then the proxy would also block requests to scss.tcd.ie. 

## Handling Requests
The proxy uses a http listener to listen for a request to come in and then will check whether the domain is on the blacklist before handling the request. If the site is not on the blacklist it will then check if the request method is CONNECT in order to choose whether to use the HTTP or HTTPS handler.

### HTTP Handler
The HTTP handler is straight forward in the sense that when the request comes in it is made by the proxy to the server and the headers and body are stored in the cache and written to the client and then the connections are closed.

### HTTPS Handler
The HTTPS handler hijacks control of the client connection then makes a connection to the to the server then simply copies the data from both. The connections are closed once all the data has been copied and transferred.
