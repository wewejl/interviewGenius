package main

import (
	"fmt"
	"interviewGenius/pkg/payment"
)

func main() {
	config := &payment.AlipayConfig{
		AppID:        "9021000134664821",
		PrivateKey:   "MIIEogIBAAKCAQEAwOjAbX3n3//o/j2WpkWULdfitzQqjPdhSbfXUNvL48FK+uVeYn13VA4OR1pZGM9alUaq+a0QJlJcz2Vt8PAcvT57JBp3dGPJc8WWAYQ/eKqwOvsF+bfa/yjrjlg4nbLjP/ZjlaM3zs39Ybnt0/vdFPDyVpNHk6H//n6Q0lKR8yL67M27RHgYkm3cU/khRvj738XJzY91vH0V2+nTMtkUorPf7WO4IRwhHGBFzhafQ/q9eenyS5tsJJ8m/eQOJ2k0dMfCjw2IK8ShT2qfBU3VD4MBKpmhdJu5IBXxdC/B+y3BimTQYC8rhwAN+MzdB/t+tOyQKbivwFZ+JkicK4Gw8wIDAQABAoIBAFGvfSQgB1rDu35EyBD6L4fF/buD/Gyap/iWPzd/CvQTOlPJYlEkPa47EXLHYCjwTLQfK3D0Bn2jrKcplQdMNW8xEOW1y1Vel8RNK3rS7CmFZYBkISCf6LzZL/2jf73PLQk9pOeNKKmKcju6hmmYIgKnEIb2cH2kQIkcQOi+jAy+y5wIxnb4+zqgbqqXzEKEPZQtV7uDIwsOKL4Nz6uokJGu5ekGe8lwiE7e7e2JRQZQAZ/CGC8XwFDaSqQVDtMFEQBS1fWR9wvQ4XlmH/IaSCy9grwoVHY+oqv5tG+tGdOJt8sKY8jJ7shWV+bcgKSnf6gOf5U7h6VvhqIvYJg5TvECgYEA7Ljb5pzqBDOmAFWZqDoJJd8z3UWHnW5JsYWtT9XIG1cfR8iFfccByFNOVUKHYLR6WAWPlYjUmRfm9aEcAaND6G79W41tfvSLZHXju5bKFbJmQqRDEOQtCELKZ2Cx0SHTBhgrk6eoMRrCIpIMraAnf3fu5VhtH943LMhiADyZ12kCgYEA0J59FcZ9vYyV3fCa6QKZsWmWQFKJZOsyhHcXII+D28hTCWhv0SIwGfgrkXgX3CAlx+CG061dxH8X5YMMkCLhIlAzOz+6Uj3y/yarThsz2SDulnRzUnwNv0u2rGhdgF+NmryoVY11ds8MbQAACQT0U9DwM8XWOPd3lxdZRszk9fsCgYBAH99pvA3kb31DT+zc1kPOH4V0JjaTXeHWleiZ3MZlKZeOoXIP3U3NT0vD6s6zUpBlsbPwhO1aP1BQL4FfrDNkDlTRbSFBJ8tuvkSfdzxs3jO3T7nfJIBSYY1krZvdk/UPDJMZX2w/SQlXxgprKhwo+nsbY3XEETUPC4UInWHrKQKBgFksxi8+r5UMuSsrpCwiDmyFw9Iu9cgLuYZiGaKzdhvGn6gP2mw8/u665HTELv7LRxsPYNKu8rwBz8cto3shTbcLLTsQXKa3EF38u5Ehk6Imr5XkpT8HBCFXTfiYjA9JyQ/xwMsBMsrcamVVcK5qTb5eO68FzDKBpb8SHfljsCNtAoGAJCdJ3lnEWm46B43kxE+3mW9D/iaJuJaFPHw/795I5XkXIyIR7RReKN/B7KofHL2RUlVoQWIaW/NU50BrSuCerWKarjp6gKNRr6C/jWz4cEFs4XYh1sCahP6QKtFgNsKOnQaLlAONkQC0HGtkGuZLVmUKSSY9VnY+1XfKBiPwyfY=",
		AliPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyX32DX+ijKSxqBCMffXjlYIKDbXPwUAs9E6g974kO3s6evhQSoMkfZYMDbqSZ9oAJQvIchE7W9sdJl3Q0prNBGhTffVlVyUDRpiKTX2WClSss8gybSRNUwY4jHZljcDFVPMlwGF/ushkuHhlg/UKuNHIeT6Xmy3nPOt7MIPeBGewDGCQlF0DYLfpha3U7puGgGMzByZ0ENVYyruCJrUMoe13Le1di5eaueElYq8G0ACVkEE8w2l+FetxN3EkaxSY1k3LffdP+0gPIm5oSPmCS3IdjCRvkn8NiwkdhyXli9ff04U7HJhQEHz01W+igIEHNpIrejpNdMFLYmsnLYpp6wIDAQAB",
		NotifyURL:    "https://your-domain.com/notify",
		ReturnURL:    "https://your-domain.com/return",
		IsProduction: false, // 测试环境
	}

	alipayService, err := payment.NewAlipayService(config)
	if err != nil {
		fmt.Printf("初始化支付宝服务失败: %v\n", err)
		return
	}

	// 创建支付订单
	payURL, err := alipayService.CreatePayment(
		"order123456", // 订单号
		"10.00",       // 金额
		"测试商品",        // 商品标题
	)
	if err != nil {
		fmt.Printf("创建支付订单失败: %v\n", err)
		return
	}

	fmt.Printf("支付链接: %s\n", payURL)
}
