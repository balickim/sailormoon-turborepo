import { Slips } from '@renderer/components/Slips'
import { createFileRoute } from '@tanstack/react-router'
import React from 'react'

export const Route = createFileRoute('/' as never)({
  component: Index
})

function Index(): React.ReactNode {
  return <Slips />
}
