import React, { useState, useEffect } from 'react';
import { ethers } from "ethers";

//const Content: React.VFC<{contract: ethers.Contract}> = ({contract}) => {
function Poll(props: { contract: ethers.Contract }) {
  const [pollTitle, setPollTitle] = useState<string>("")
  const [pollResult, setPollResult] = useState<number>(0)
  const [message, setMessage] = useState<string>("")
  const { getPoll, checkPoll, vote } = props.contract.functions;
  useEffect(() => {
    const getPollFunc = async () => {
      setPollTitle(await getPoll())
      const result = (await checkPoll()).toString()
      setPollResult(result)
      switch (result) {
        case "0":
          setMessage("push agree or disagree")
          break
        case "1":
          setMessage("you had agreed")
          break
        case "2":
          setMessage("you had disagreed")
          break
      }
    }
    getPollFunc()
  }, [pollResult])

  const hundleSubmit = async (poll: number) => {
    if (pollResult > 0) return
    await vote(poll)
    setPollResult(poll)
  }

  return (
    <>
      <div>{pollTitle}</div>
      <div>
        <button onClick={() => { hundleSubmit(1) }}>agree</button>
        <button onClick={() => { hundleSubmit(2) }}>disagree</button>
      </div>
      <div>{message}</div>
    </>
  )
}

export default Poll;