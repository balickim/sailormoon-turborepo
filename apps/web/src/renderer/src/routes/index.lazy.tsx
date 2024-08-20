import { createLazyFileRoute } from '@tanstack/react-router'
import React from 'react'

export const Route = createLazyFileRoute('/' as never)({
  component: Index
})

function Index(): React.ReactNode {
  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
    </div>
  )
}
