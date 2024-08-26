import { Link, Outlet } from '@tanstack/react-router'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout')({
  component: () => (
    <>
      <div className="navbar bg-base-100">
        <div className="navbar-start"></div>

        <div className="navbar-center hidden lg:flex">
          <ul className="menu menu-horizontal px-1">
            <li>
              <Link to="/" className="[&.active]:font-bold">
                Keje
              </Link>
            </li>
            <li>
              <Link to="/boats" className="[&.active]:font-bold">
                Łodzie
              </Link>
            </li>
          </ul>
        </div>
      </div>

      <hr />
      <Outlet />
    </>
  )
})
