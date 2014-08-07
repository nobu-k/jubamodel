jubamodel
=========

`jubamodel` is a tool to display information about machine learning models
created by the [Jubatus](http://jubat.us) software and convert them for
compatibility with various Jubatus versions.

Installation
------------

`go get github.com/nobu-k/jubamodel`

Usage
-----

* To display model information:  
  `jubamodel info file [files...]`
* To convert models:  
  `jubamodel rewrite-version file new-version`
