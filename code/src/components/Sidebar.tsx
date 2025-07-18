"use client";

import { type ReactElement, type ReactNode} from "react";
import MobileSidebar from "./MobileSidebar";
import { usePathname } from 'next/navigation';

interface MenuItem {
    name: string;
    icon: ReactElement;
    href: string;
}

interface SidebarProps {
  menuItems: MenuItem[];
  accountDropdown?: ReactNode;
}

export default function Sidebar({ menuItems, accountDropdown }: SidebarProps) {
  const pathname = usePathname();
  
  // Find the menu item that matches the current path
  const activeItem = menuItems.find(item => pathname.startsWith(item.href))?.name || menuItems[0]?.name;

  return (
    <>
      {/* Mobile Hamburger Button */}
      <MobileSidebar 
        menuItems={menuItems}
        accountDropdown={accountDropdown}
      />

      {/* Desktop Sidebar */}
      <aside
        className="hidden lg:block fixed inset-y-0 left-0 w-64 bg-neutral-800 border-neutral-700 z-50"
        aria-label="Sidebar"
      >
        <div className="relative flex flex-col h-full">
          <header className="p-4 mt-1 ml-1 flex justify-between items-center gap-x-2">
            <a
              className="flex-none font-semibold text-xl"
              href="#"
              aria-label="Brand"
            >
              Ekstra
            </a>
          </header>

          <nav className="flex-1 overflow-y-auto">
            <div className="p-2 w-full flex flex-col flex-wrap">
              <ul className="space-y-1">
                {menuItems.map((item) => (
                  <li key={item.name}>
                    <a
                      className={`w-full flex items-center gap-x-3.5 p-2.5 text-base rounded-lg hover:bg-neutral-700 ${
                        activeItem === item.name ? "bg-neutral-700" : ""
                      }`}
                      href={item.href}
                    >
                      {item.icon}
                      {item.name}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </nav>

          {/* Account Info Section */}
          <footer className="p-4 border-t border-neutral-700">
            {accountDropdown}
          </footer>
        </div>
      </aside>
    </>
  );
}