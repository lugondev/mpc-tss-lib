
?
mpc_client.protopb"
EmptyParams"

ReplyEmpty"(
GetPartyParams
pubkey (	Rpubkey"T
GetPartyResponse
id (	Rid
address (	Raddress
pubkey (	Rpubkey"X

PartyShare
pubkey (	Rpubkey
address (	Raddress
partyId (	RpartyId"
GetPartiesParams"<
GetPartiesResponse&
shares (2.pb.PartyShareRshares"h

SignParams
id (	Rid
parties (	Rparties
message (Rmessage
pubkey (	Rpubkey"V
SignResponse
id (	Rid
message (	Rmessage
	signature (	R	signature"9
KeygenGeneratorParams
id (	Rid
ids (	Rids"1
KeygenGeneratorResponse
pubkey (	Rpubkey"&
RequestPartyResponse
id (	Rid" 
Pong
message (	Rmessage2?
MpcParty;
RequestParty.pb.EmptyParams.pb.RequestPartyResponse" K
KeygenGenerator.pb.KeygenGeneratorParams.pb.KeygenGeneratorResponse" 6
GetParty.pb.GetPartyParams.pb.GetPartyResponse" <

GetParties.pb.GetPartiesParams.pb.GetPartiesResponse" *
Sign.pb.SignParams.pb.SignResponse" #
Ping.pb.EmptyParams.pb.Pong" BZpb/bproto3