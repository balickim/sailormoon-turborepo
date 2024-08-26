import React, { forwardRef, useImperativeHandle, useState } from 'react'
import { Description, Dialog, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react'

interface IDialogBaseProps {
  title: string
  children: React.ReactNode
  ok?: string
  cancel?: string
}

const DialogBase = forwardRef((props: IDialogBaseProps, ref) => {
  const [isOpen, setIsOpen] = useState(false)

  useImperativeHandle(ref, () => ({
    open: () => setIsOpen(true),
    close: () => setIsOpen(false)
  }))

  return (
    <Dialog open={isOpen} onClose={() => setIsOpen(false)} className="relative z-50">
      <DialogBackdrop className="fixed inset-0 backdrop-blur-sm" />

      <div className="fixed inset-0 flex w-screen items-center justify-center p-4">
        <DialogPanel className="max-w-lg space-y-4 border bg-white p-12">
          <DialogTitle className="font-bold">{props.title}</DialogTitle>
          <Description>{props.children}</Description>
          {props.cancel || props.ok ? (
            <div className="flex gap-4">
              {props.cancel ? <button onClick={() => setIsOpen(false)}>Cancel</button> : null}
              {props.ok ? <button onClick={() => setIsOpen(false)}>Deactivate</button> : null}
            </div>
          ) : null}
        </DialogPanel>
      </div>
    </Dialog>
  )
})
DialogBase.displayName = 'DialogBase'

export default DialogBase
