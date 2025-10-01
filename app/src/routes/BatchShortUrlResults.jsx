import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import styled from 'styled-components'

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
`

const Table = styled.table`
  width: 100%;
  max-width: 800px;
  border-collapse: collapse;
  margin-top: 20px;
`

const Th = styled.th`
  border: 1px solid #ddd;
  padding: 12px;
  text-align: left;
  background: #201f30;
  color: white;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-radius: 4px 4px 0 0;
`

const Td = styled.td`
  border: 1px solid #ddd;
  padding: 12px;
`

const ShortUrlLink = styled.a`
  color: #007bff;
  text-decoration: none;
  &:hover {
    text-decoration: underline;
  }
`

const BatchShortUrlResults = () => {
  const { batch_id } = useParams()
  const [urlsMap, setUrlsMap] = useState([])
  const appUrl = window.location.origin

  const fetchShortUrlsByBatchId = async () => {
    if (!batch_id) return

    try {
      const res = await fetch(`/api/all/${batch_id}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      })
      const data = await res.json()
      console.log('debug res', data?.data?.urls_map)
      setUrlsMap(data?.data?.urls_map)
    } catch {
      alert('Unexpected error occurred. Make sure your batch id is correct and not expired.')
      setUrlsMap([])
    }
  }

  useEffect(() => {
    if (!batch_id) return

    fetchShortUrlsByBatchId()
  }, [batch_id])

  return (
    <Container>
      <h1>Batch Short URL Results for {batch_id}</h1>
      {urlsMap && (
        <Table>
          <thead>
            <tr>
              <Th>Long Url</Th>
              <Th>Short Url</Th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(urlsMap).map(([longUrl, shortUrl], index) => (
              <tr key={index}>
                <Td>{longUrl}</Td>
                <Td>
                  <ShortUrlLink href={`${appUrl}/${shortUrl}`} target="_blank" rel="noopener noreferrer">
                    {`${appUrl}/${shortUrl}`}
                  </ShortUrlLink>
                </Td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </Container>
  )
}

export default BatchShortUrlResults