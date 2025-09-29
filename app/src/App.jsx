import { Routes, Route } from 'react-router-dom'
import styled from 'styled-components'
import './App.css'
import Home from './routes/Home'
import RedirectPage from './routes/RedirectPage'

const Hero = styled.div`
  position: relative;
  width: 100%;
  height: 100dvh; /* avoid 100vw/100vh scrollbar issues */
  display: flex;
  flex-direction: column;
  align-items: center; /* center children horizontally */
  box-sizing: border-box; /* include padding in the element's height */
  padding-block: 32px; /* vertical padding */
`

function App() {
  return (
    <Hero>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/:short_code" element={<RedirectPage />} />
      </Routes>
    </Hero>
  )
}
export default App
