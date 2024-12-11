<div align="center">

  <h1><a href="https://github.com/Karthik-HR0/refx">refx</a></h1>

> refx is an automated tool for finding reflected parameters on web applications!

</div>

<hr>
<br>

## `features`

- Fast and efficient crawling of domains
- Real-time detection of reflected parameters
- Crawls subdomains
- Saves crawled pages and parameters to structured folders
- Easy integration with other security testing tools
- Displays results in the terminal and saves to a file

<br>
<br>

`installation`

```bash
go install github.com/Karthik-HR0/refx@latest
```
---
## Usage 
<br>
<br><pre>
options:
  -h, -help      show help message
  -t             target URL (e.g., http://example.com)
  -s             crawl subdomains
</pre><br>

---

## Example 

```bash
 refx -t http://example.com
```
---
```
./refx -t http://testphp.vulnweb.com      

 
 
              ___________       
_______   ____\_   _____/__  ___
\_  __ \_/ __ \|    __) \  \/  /
 |  | \/\  ___/|     \   >    <     
 |__|    \___  >___  /  /__/\_ \
             \/    \/         \/    
                                    V1.0
                                    @Karthik-HR0
                                    
    Automated Reflected Parameter Finder Tool
                       
  
[*] Crawling the domain for pages and parameters...
[*] Crawled 42 unique pages
[*] Found 5 unique parameters
[*] Testing for reflected parameters...

[+] Reflected Parameters Found:
[Reflected] http://testphp.vulnweb.com/listproducts.php?artist=reflect_test_parameter
[Reflected] http://testphp.vulnweb.com/listproducts.php?cat=reflect_test_parameter
[Reflected] http://testphp.vulnweb.com/showimage.php?file=reflect_test_parameter
```

---
## troubleshooting

Ensure the target is accessible

Check your internet connection

Verify URL format (http:// or https://)

use -t for target 


<br>
<br>
<br>
<p align="center">
Made with <3 by <a href="https://github.com/Karthik-HR0">@Karthik-HR0</a>
<br>

</p>
