import CreatorAccountDropdown from "~/components/CreatorAccountDropdown";
import { menuItems } from "~/components/CreatorMenuItems";
import { retrieveCreatorDropdownData } from "~/server/actions";
import Sidebar from "~/components/Sidebar";

export default async function creatorLayout({
  children,
}: {
  children: React.ReactNode;
}) {

  const creatorData = await retrieveCreatorDropdownData();

  return(
    <>
      <div className="lg:ml-64">
        <Sidebar 
          menuItems={menuItems}
          accountDropdown={ <CreatorAccountDropdown creatorDropdownData={creatorData}/>}
        />
        {children}
      </div>
    </>
  )
}