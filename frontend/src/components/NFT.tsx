import React, { useState, useEffect, useRef } from 'react';
import { ethers } from "ethers";
import artifact from "../abi/MyNFT.json"
import axios from "axios"

function NFT(props: { signer_address: string, friend_address: string }) {
  type TokenItem = {
    tokenID: string
    tokenURI: string
    tokenInfo: TokenInfo
  }
  type TokenInfo = {
    name: string
    description: string
    image: string
    attributes: string
  }
  const provider = new ethers.providers.JsonRpcProvider();
  const nft_address = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
  const signer = provider.getSigner(props.signer_address);
  const contract = new ethers.Contract(nft_address, artifact.abi, provider);
  const client = contract.connect(signer);

  const inputElement = useRef<HTMLInputElement>(null)

  const [balance, setBalance] = useState<string>("")
  const [isOwner, setIsOwner] = useState<boolean>(false)
  const [tokenName, setTokenName] = useState<string>("")
  const [tokenItems, setTokenItems] = useState<Array<TokenItem>>([])
  const { name, tokenURI, balanceOf, safeMint, approve, transferFrom } = client.functions;

  useEffect(() => {
    const setTokenNameFunc = async () => {
      setTokenName((await name()).toString())
    }
    const setIsOwnerFunc = async () => {
      const owner = (await client.owner()).toString()
      const signer = (await client.signer.getAddress()).toString()
      setIsOwner(owner === signer)
    }

    setTokenNameFunc()
    setIsOwnerFunc()
    setBalanceFunc()
    setTokenItemsFunc()
  }, [])

  const setBalanceFunc = async () => {
    const ret = (await balanceOf(props.signer_address)).toString()
    setBalance(await ret)
  }

  const setTokenItemsFunc = async () => {
    let items: Array<TokenItem> = [];
    const sentLogs = await client.queryFilter(
      client.filters.Transfer(props.signer_address, null),
    );
    const receivedLogs = await client.queryFilter(
      client.filters.Transfer(null, props.signer_address),
    );

    const logs = sentLogs.concat(receivedLogs).sort((a, b) =>
      a.blockNumber - b.blockNumber ||
      a.transactionIndex - b.transactionIndex,
    )

    const owned = new Set<string>()
    for (const log of logs) {
      if (log?.args?.to === (await client.signer.getAddress()).toString()) {
        owned.add(log?.args?.tokenId.toString())
      } else {
        owned.delete(log?.args?.tokenId.toString())
      }
    }

    for (const id of owned) {
      const uri = (await tokenURI(id)).toString()
      const token_info = await axios.get(uri)
      console.log(token_info)
      const item: TokenItem = {
        tokenID: id,
        tokenURI: uri,
        tokenInfo: token_info.data as unknown as TokenInfo,
      }
      items.push(item)
    }
    setTokenItems(items)
  }
  
  const hundleMint = async () => {
    safeMint(props.signer_address, "http://localhost:3000/nft.json").then(() => {
      setBalanceFunc()
      setTokenItemsFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const hundleApprove = async(tokenID: string) => {
    await approve(props.friend_address, tokenID)
  }

  const hundleTransferToMe = async () => {
    const id = inputElement?.current?.value
    transferFrom(props.friend_address, props.signer_address, id).then(() => {
      setBalanceFunc()
      setTokenItemsFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const hundleTransfer = async(id: string) => {
    transferFrom(props.signer_address, props.friend_address, id).then(() => {
      setBalanceFunc()
      setTokenItemsFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const tokenItem = (item: TokenItem) => {
    return(
      <div key={item.tokenID} style={{ margin: "10px 0" }}>
        <img src={item.tokenInfo.image} style={{width: "200px"}} />
        <div>id: {item.tokenID}</div>
        <div>name: {item.tokenInfo.name}</div>
        <div>description: {item.tokenInfo.description}</div>
        <div>
          <button onClick={() => { hundleTransfer(item.tokenID) }}>transfer to friend</button>
          <button onClick={() => { hundleApprove(item.tokenID) }}>approve to transfer from friend</button>
        </div>
      </div>
    )
  }

  return (
    <>
    <div style={{margin: "10px 20px", borderRadius: "10px", boxShadow: "2px 2px 10px #bbb", padding: "20px"}}>
      <div>{tokenName}</div>
      <div>user_address: {props.signer_address}{isOwner ? " (owner)" : ""}</div>
      <div><input ref={inputElement} placeholder="input nft id" /><button onClick={() => { hundleTransferToMe() }}>transfer to me</button></div>
      {isOwner ?
        <div>
          <button onClick={() => { hundleMint() }}>mint</button>
        </div> :
        null
      }
      <div>nfts: {balance}items</div>
      {tokenItems.map(tokenItem)}
    </div>
    </>
  )
}

export default NFT;
