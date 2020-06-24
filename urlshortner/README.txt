URL Shortner README

	Thank you for using my URL Shortner. For this URL shortner, it is only in an API test state, meaning an API tester will be needed for this (examples include Talend API Chrome Extension, Postman, and others that can do the same job). Once the program is ran, you will know as in the terminal, a message called "Running on port 8091" will be displayed.
	Once running successfully, use the API test, make a PUT request and enter "http://localhost:8091/create" as the request URL. Go to the body section of data to enter, make sure it is set to JSON. After that, make sure it is formatted correctly to look what it is below and enter the URl you wish to shorten in the indicated area. 

{
	"longUrl": "Enter the URL here" 
}

Once that is done, send the request, and you can observe the response below. Essentially, another JSON file gets sent back with 3 fields of data provided
{ID: 7fh62sd (Will be random), longUrl: the original URL entered/URL desired to be shortened, shortUrl: http://localhost:8091/7fh62sd}

NOTICE, the ID is attached at the end of the shortUrl

You can simply copy and paste the shortUrl into a browser and into the search URL and should redirect you to the original longUrl you desired to shorten.


