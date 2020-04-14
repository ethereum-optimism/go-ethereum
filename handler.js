function handleRequest(req, res) {
  if(req.method == "eth_getBalance") {
    res.result = "0x99"
  }

  return res
}
