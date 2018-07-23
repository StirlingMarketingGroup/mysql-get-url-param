
# MySQL `get_url_param`
 
This is the first MySQL UDF I've seen that can actually user arguments and that is written in Go. So besides this  UDF's specific task of parsing URL parameters, this will make an excellent example for making new UDF's with Go in the future.
 
Go is exceptionally powerful and fairly simple language to understand and write in, with a good community behind it and a strong built in library. It also can be used directly with C programs, which it made it sound perfect for creating MySQL UDF's with way more ease than using pure C. 

I mean here I am, a lowley PHP dev, and now my MySQL instances have the power to extract URL parameters with perfection. 

## Installation instructions

These are the Linux instructions, so but I'm sure the standard, say Windows, instructions can be adapted to work with the Go steps of compiling and installing the plugin.

1. make sure libmysqlclient-dev is installed
    `apt install libmysqlclient-dev`
2. Find your MySQL plugins dir, which is likely "/usr/lib/mysql/plugin" (or you can find it by running `select @@plugin_dir;` on your MySQL server)
3. Navigate to the folder where you cloned this repository to, and run the following (replacing the plugin path with your own)
    `go build -buildmode=c-shared -o get_url_param.so && cp get_url_param.so /usr/lib/mysql/plugin/get_url_param.so`
4. Then on your MySQL server, run this to declare the function in your current schema
    ``create function`get_url_param`returns string soname'get_url_param.so';``

## Usage
> string  **get_url_param** ( string  `url` , string  `parameter` )

### Parameters
**url** - The URL containing the query string (URL doesn't need to be complete, e.g. can start with "http://...", or "/...", etc.)
**parameter** - The name of the parameter to have its value returned

### Return values
This function will **always** return a string if the parameter exists, even its value is empty. Otherwise, if the parameter does not exist, this function will return `null`.

If a string is returned, it is URL decoded, including plus signs (+), so `%20` and `+` in a parameter both will be converted to ` `. This also fully supports UTF8 (and likely many other character sets, considering it's using the built-in Go URL parser).

### Examples

```mysql
select`get_url_param`('https://www.youtube.com/watch?v=KDszSrddGBc','v');
-- "KDszSrddGBc"
```

```mysql
select`get_url_param`('watch?v=KDszSrddGBc','v');
-- "KDszSrddGBc"
```

```mysql
select`get_url_param`('watch?v=KDszSrddGBc','x');
-- null
```

```mysql
select`get_url_param`('https://www.google.com/search?q=cgo+uint32+to+pointer&rlz=1C1CHBF_enUS767US767&oq=cgo+uint32+to+pointer&aqs=chrome..69i57.12106j0j7&sourceid=chrome&ie=UTF-8','q');
-- "cgo uint32 to pointer"
```

```mysql
select`get_url_param`('/search?q=Na%C3%AFvet%C3%A9&oq=Na%C3%AFvet%C3%A9','q');
-- "Naïveté"
```
