import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom'

const RedirectPage = () => {
  const { short_code } = useParams()

  useEffect(() => {
    if (!short_code) return
    
    window.location.replace(`/api/${short_code}`)
  }, [short_code])

  return <div>Redirecting…</div>
}

export default RedirectPage