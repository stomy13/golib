import { useState } from 'react'
import './App.css'
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

// Import service definition that you want to connect to.
import { PingService } from "./gen/connect/ping/v1/ping_pb";

// The transport defines what type of endpoint we're hitting.
// In our example we'll be communicating with a Connect endpoint.
// If your endpoint only supports gRPC-web, make sure to use
// `createGrpcWebTransport` instead.
const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

// Here we make the client itself, combining the service
// definition with the transport.
const client = createClient(PingService, transport);

function App() {
  const [inputValue, setInputValue] = useState("");
  const [messages, setMessages] = useState<
    {
      fromMe: boolean;
      message: string;
    }[]
  >([]);
  return <>
      <ol>
      {messages.map((msg, index) => (
        <li key={index}>
          {`${msg.fromMe ? "ME:" : "ELIZA:"} ${msg.message}`}
        </li>
      ))}
    </ol>
    <form onSubmit={async (e) => {
      e.preventDefault();
      // Clear inputValue since the user has submitted.
      setInputValue("");
      // Store the inputValue in the chain of messages and
      // mark this message as coming from "me"
      setMessages((prev) => [
        ...prev,
        {
          fromMe: true,
          message: inputValue,
        },
      ]);
      const response = await client.ping({
        number: BigInt(123456),
        text: inputValue,
      });
      setMessages((prev) => [
        ...prev,
        {
          fromMe: false,
          message: `number: ${response.number}, text: ${response.text}`,
        },
      ]);
    }}>
      <input value={inputValue} onChange={e => setInputValue(e.target.value)} />
      <button type="submit">Send</button>
    </form>
    </>
}

export default App
