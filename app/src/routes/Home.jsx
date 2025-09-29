import React, { useState } from 'react'
import styled from 'styled-components'

const SM_VIEW = '620px'

const Container = styled.div`
  position: relative; /* establish containing block for absolute child */
  display: flex;
  flex-direction: column;
  align-items: center;      /* horizontal center for heading */
  width: 100%;
  height: 100%;
  flex: 1; /* fill the Hero container */
  text-align: center;
  gap: 16px;
`

const UrlContainer = styled.div`
  border: 1px solid rgba(0,0,0,0.10);
  box-shadow: 0 8px 20px rgba(0,0,0,0.15);
  box-sizing: border-box; /* include padding in width calc */
  width: 100%;          /* take available width of parent */
  max-width: 800px;     /* cap at 700px on larger screens */
  min-width: 400px;         /* allow shrinking on small screens */
  padding: 32px 16px;
  margin-top: 128px;
  margin-bottom: 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
`

const UrlInputContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
  padding: 0 8px;
  box-sizing: border-box; /* include padding in width calc */

  /* On very small viewports, allow shrinking below 400px to avoid overflow */
  @media (max-width: ${SM_VIEW}) {
    flex-direction: column;
    min-width: 0;
    max-width: 100%;
  }
`

const UrlInput = styled.input`
  flex: 1 1 auto;   /* allow input to shrink and grow inside the container */
  min-width: 0;     /* critical for flex children to shrink below content size */
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  line-height: 1.2;
  border: 1px solid #d1d5db; /* gray-300 */
  border-radius: 10px;
  outline: none;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  margin-right: 8px;

  /* Remove UA outline on mouse focus; keep visible ring for keyboard */
  &:focus { outline: none; }

  &:focus-visible {
    border-color: #f59e0b; /* amber-500 */
    box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.35);
  }

  @media (max-width: ${SM_VIEW}) {
    margin-right: 0;      /* no horizontal gap in column layout */
    margin-bottom: 8px;   /* vertical gap above button */
  }
`;

const ShortenBtn = styled.button`
  background: linear-gradient(135deg, #f59e0b, #ea580c); /* amber → orange */
  color: white;
  border: none;
  padding: 16px 16px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 10px; /* match input radius */
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.05s ease;
  white-space: nowrap;
  width: 150px;
  outline: none;
  -webkit-tap-highlight-color: transparent; /* remove tap highlight on mobile */

  &:hover {
    transform: translateY(-1px);
    background: linear-gradient(135deg, #d97706,rgb(191, 68, 19)); /* darker on press */
    box-shadow: 0 10px 20px rgba(234, 95, 20, 0.35); /* orange shadow */
  }

  &:active {
    transform: translateY(0);
    background: linear-gradient(135deg, #d97706, #c2410c); /* darker on press */
    box-shadow: 0 4px 10px rgba(194, 65, 12, 0.3);
    outline: none;
  }

  /* Remove default UA focus ring on mouse focus; keep focus-visible for keyboard */
  &:focus {
    outline: none;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  }

  &:focus-visible {
    outline: none;
    box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.35), 0 1px 2px rgba(0, 0, 0, 0.06);
  }

  /* Firefox inner focus reset */
  &::-moz-focus-inner {
    border: 0;
  }

  ${({ disabled }) => disabled && `
    background:rgba(52, 52, 52, 0.81);
    color: grey;
    cursor: not-allowed;
    pointer-events: none;
  `}

`

/* Result section styles */
const ResultCard = styled.div`
  margin-top: 24px;
  padding: 20px 24px;
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.10);
  box-shadow: 0 8px 20px rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  max-width: 640px;
  width: 100%;
`

const ResultTitle = styled.h4`
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.2px;
  color: white;
`

const ResultLink = styled.p`
  padding: 8px 16px;
  cursor: pointer;
  border-radius: 999px;
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: #ffffff;
  text-decoration: none;
  font-weight: 700;
  letter-spacing: 0.2px;
  box-shadow: 0 6px 14px rgba(22,163,74,0.35);
  transition: transform 0.12s ease, box-shadow 0.12s ease, opacity 0.2s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 10px 20px rgba(22,163,74,0.35);
    background: linear-gradient(135deg,rgb(25, 136, 66),rgb(17, 132, 59));
  }

  &:active {
    transform: translateY(0);
    background: linear-gradient(135deg,rgb(19, 107, 51),rgb(14, 111, 49));
  }
`

const Home = () => {
  const [longUrl, setLongUrl] = useState('')
  const [shortCode, setShortCode] = useState('')
  const [error, setError] = useState(null)
  const [copied, setCopied] = useState(false)
  const appUrl = window.location.origin

  const handleShorten = async () => {
    if (!longUrl) return

    setError(null)
    setShortCode(null)
    try {
      const res = await fetch(`/api/shorten`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ long_url: longUrl }),
      })
      const data = await res.json()
      const shortUrl = data?.data?.short_url
      setShortCode(shortUrl)
    } catch (err) {
      setError('Ops...Something went wrong.')
      setShortCode(null)
      console.error('Shorten request failed:', err)
    } finally {
      setLongUrl('')
    }
  }

  const handleCopyShortCode = async () => {
    if (!shortCode) return

    try {
      await navigator.clipboard.writeText(`${appUrl}/${shortCode}`)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch (e) {
      console.error('Failed to copy', e)
    }
  }

  return (
    <Container>
      <h1>Welcome to Short URL</h1>
      <UrlContainer>
        <h3>Shorten your long URL</h3>
        <UrlInputContainer>
          <UrlInput
            type="url"
            placeholder="Enter a long URL..."
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            required
          />
          <ShortenBtn disabled={!longUrl} onClick={handleShorten}>Shorten URL</ShortenBtn>
        </UrlInputContainer>
      </UrlContainer>
      {shortCode && (
        <ResultCard>
          <ResultTitle>Here is your shortened URL</ResultTitle>
          <ResultLink
            role="button"
            tabIndex={0}
            onClick={handleCopyShortCode}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                handleCopyShortCode()
              }
            }}
            aria-label="Copy short code to clipboard"
            title="Click to copy code"
          >
            {`${appUrl}/${shortCode}`}
          </ResultLink>
          {copied && <p style={{ color: 'orange' }}>Copied to clipboard</p>}
        </ResultCard>
      )}
      {error && (
        <ResultCard>
          {error}
        </ResultCard>
      )}
    </Container>
  )
}

export default Home