import { forwardRef } from 'react'
import DialogBase from '../DialogBase'
import { Formik, Form, Field, ErrorMessage } from 'formik'
import { z } from 'zod'
import { toFormikValidationSchema } from 'zod-formik-adapter'

const slipSchema = z.object({
  slipName: z.string().min(1, 'Slip name is required'),
  slipAmount: z.number().min(1, 'Amount must be greater than 0')
})

const AddSlipDialog = forwardRef((props, ref) => {
  return (
    <DialogBase title="Dodawanie" ref={ref}>
      <Formik
        initialValues={{
          slipName: '',
          slipAmount: ''
        }}
        validationSchema={toFormikValidationSchema(slipSchema)}
        onSubmit={(values, { setSubmitting }) => {
          console.log('Form Submitted:', values)
          setSubmitting(false)
          if (ref.current) {
            ref.current.close()
          }
        }}
      >
        {({ isSubmitting }) => (
          <Form>
            <div className="mb-4">
              <label htmlFor="slipName" className="block text-sm font-medium text-gray-700">
                Slip Name
              </label>
              <Field
                type="text"
                name="slipName"
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
              />
              <ErrorMessage name="slipName" component="div" className="text-red-500 text-sm mt-1" />
            </div>

            <div className="mb-4">
              <label htmlFor="slipAmount" className="block text-sm font-medium text-gray-700">
                Slip Amount
              </label>
              <Field
                type="number"
                name="slipAmount"
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
              />
              <ErrorMessage
                name="slipAmount"
                component="div"
                className="text-red-500 text-sm mt-1"
              />
            </div>

            <div className="mt-4">
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-4 py-2 bg-blue-500 text-white rounded"
              >
                Submit
              </button>
            </div>
          </Form>
        )}
      </Formik>
    </DialogBase>
  )
})

AddSlipDialog.displayName = 'AddSlipDialog'

export default AddSlipDialog
