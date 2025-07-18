import AccountDropdown from "~/components/AccountDropdown";
import Sidebar from "~/components/Sidebar";
import { menuItems } from "~/components/UserMenuItems";

export default function userLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return(
    <>
      <div className="lg:ml-64">
        <Sidebar 
          menuItems={menuItems}
          accountDropdown={ <AccountDropdown />}
        />
        {children}
      </div>
    </>
  )
}