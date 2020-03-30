package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"go.zoe.im/unipaw"
	"go.zoe.im/x"
	"go.zoe.im/x/cli"

	. "go.zoe.im/unipaw"
)

var (
	tokenRegex = regexp.MustCompile("Auth\\=(\\w+)")
)

func turnJSON(body []byte) (interface{}, error) {
	var d map[string]interface{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func bodyHookWei(body []byte) (interface{}, error) {
	var d map[string]interface{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	ds, ok := d["degrees"]
	if !ok {
		return nil, errors.New("数据不符合预期,应该存在 degrees 字段")
	}

	dss, ok := ds.(map[string]interface{})
	if !ok {
		return nil, errors.New("数据不符合预期,应该存在 degrees.degree 字段")
	}
	return dss["degree"], nil
}

func bodyHookLi(body []byte) (interface{}, error) {
	var d map[string]interface{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	if d["result"] != nil && d["result"].(string) == "false" {
		return nil, fmt.Errorf("%s", d["message"])
	}

	return d, nil
}

func init() {

	// create the site
	var site = NewSite(
		SiteName("zwfw"),
		SiteVersion("v0"),
		SiteHost("http://app.gjzwfw.gov.cn"),
		SiteURL("jimp/jiaoyubu/interfaces"),
		// add path as a option
		SitePath(
			PathName("xuewei"),
			PathURL("xuewei.do"),
			PathBodyHook(bodyHookWei),
			PathParam(
				ParamName("xm"),
			),
			PathParam(
				ParamName("zsbh"),
				ParamQuery("xwzsbh"),
			),
		),
	)

	// register a path and get the return
	var tokenPath = site.RegisterWithOptions(
		PathName("token"),
		PathURL("yzticket.do"),
		PathBodyHook(func(body []byte) (interface{}, error) {
			res := tokenRegex.FindSubmatch(body)
			if len(res) != 2 {
				return "", errors.New("find token from body error")
			}
			return x.Bytes2Str(res[1]), nil
		}),
	)

	// register a path
	site.RegisterWithOptions(
		PathName("xueli"),
		PathURL("xueli.do"),
		PathBodyHook(bodyHookLi),
		PathParam(
			ParamName("xm"),
		),
		PathParam(
			ParamName("zsbh"),
		),
		PathParam(
			ParamName("token"),
			ParamValue(func(r *http.Request) (interface{}, error) {
				return tokenPath.Do()
			}),
		),
	)

	// register the site
	Register(site)
}

func main() {
	var addr string
	var prefix string
	cli.New(
		cli.Name("unipaw"),
		cli.Short("A univesal paw platform to generate api for every site."),
		cli.SetFlags(func(c *cli.Command) {
			c.Flags().StringVarP(&addr, "addr", "", ":8080", "address to listen.")
			c.Flags().StringVarP(&prefix, "prefix", "", "/api", "endpoint prefix.")
		}),
		cli.Run(func(c *cli.Command, args ...string) {
			unipaw.Init(unipaw.Prefix(prefix))
			log.Printf("serve http %s\n", addr)
			err := http.ListenAndServe(addr, unipaw.DefaultManager)
			if err != nil {
				log.Printf("exit with error %s\n", err)
			}
		}),
	).Run()
}
