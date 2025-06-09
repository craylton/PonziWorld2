import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

// Types for state
interface HelloResponse {
  message: string
}

function App() {
  const [count, setCount] = useState<number>(0)
  const [message, setMessage] = useState<string>('')

  // Fetch hello message from backend
  useEffect(() => {
    const fetchMessage = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/hello')
        const data: HelloResponse = await res.json()
        setMessage(data.message)
      } catch {
        setMessage('Error fetching from backend.')
      }
    }
    fetchMessage()
  }, [])

  return (
    <>
      {/* Logos */}
      <div>
        <a href="https://vite.dev" target="_blank" rel="noopener noreferrer">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank" rel="noopener noreferrer">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
      <h1>PonziWorld2</h1>
      <p>{message}</p>
    </>
  )
}

export default App
