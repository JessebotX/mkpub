# pub

- Easy to use
- Static site generator that also supports creating RSS feeds, EPUBs, etc. from your Markdown files (WIP)
- Highly customizable and extensible: customize the underlying layout using `go`'s advanced templating engine
- Great defaults; zero/minimal configuration required
- No lock-in to any service or ecosystem
- Local and private
- Backward and forward-compatible (at and after version 1.0.0)

## Description

`pub` is an easy-to-use *static site generator* specifically **designed for independently publishing your Markdown-based written works on the web**. `pub` also supports *generating other popular open formats such as PDFs and EPUBs (WIP)*, which are supported across almost all machines and e-readers.

*`pub` requires minimal configuration*, allowing you to put all your focus on the act of writing itself. These great defaults emphasize user accessibility and widespread compatibility with many systems, letting you produce works that are comfortable to read in with little effort. Nevertheless, *`pub` is also extremely configurable and extensible*, with an embedded text templating engine (`go`'s `text/template` and `html/template`) and support for custom configuration fields that allows you to develop and utilize custom themes to modify the overall presentation of the output.

Additionally, `pub` is permissively licensed. You can use the program and browse the documentation locally without internet, and with all your writing being done in open plain-text formats, there is little potential lock-in to any single service/ecosystem.

The upcoming `pub` version 1.0.0 is also committed to not only be backward-compatible but also forward-compatible. Old book projects should be able to be successfully built in the future.

## Usage
