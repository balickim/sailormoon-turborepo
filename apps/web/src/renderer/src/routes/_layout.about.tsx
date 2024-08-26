import { createFileRoute } from '@tanstack/react-router'
import React from 'react'

export const Route = createFileRoute('/about' as never)({
  component: Index
})

function Index(): React.ReactNode {
  return <>About</>
}
