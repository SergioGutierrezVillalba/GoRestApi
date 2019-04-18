package jwt

import ()

type JWTHeader struct {
	AlgorithmUsed		string
	TypeOfToken			string
}

func NewHeader() JWTHeader{
	header := JWTHeader{
		AlgorithmUsed: "",
		TypeOfToken: "JWT",
	}
	return header
}

func NewHeaderFromRaw() JWTHeader{
	header := JWTHeader{
		AlgorithmUsed: "",
		TypeOfToken: "JWT",
	}
	return header
}