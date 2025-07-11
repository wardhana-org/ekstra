'use client';

import { useState, type SetStateAction } from 'react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline';
import { menuItems } from '~/components/UserMenuItems';
import AccountDropdown from './AccountDropdown';

export default function MobileSidebarHome() {
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);
  const [activeItem, setActiveItem] = useState("Home");

    const handleItemClick = (itemName: SetStateAction<string>) => {
    setActiveItem(itemName);
    setIsMobileSidebarOpen(false);
  };

  return (
    <>
      {/* Mobile Hamburger Button */}
      <button
        onClick={() => setIsMobileSidebarOpen(true)}
        className="lg:hidden fixed top-4 right-4 z-50 p-2 text-white 
                  transition-all duration-300"
        aria-label="Open sidebar"
      >
        <Bars3Icon className="h-6 w-6 scale-150 cursor-pointer" />
      </button>

      {/* Mobile Sidebar - Fullscreen */}
      {isMobileSidebarOpen && (
        <div className="lg:hidden fixed inset-0 z-50">
          <aside
            className={`absolute inset-0 bg-white dark:bg-neutral-800 transform ${
              isMobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'
            } transition-transform duration-300 ease-in-out`}
          >
            <div className="relative flex flex-col h-full max-h-full">
              <header className="p-4 mt-1 ml-1 flex justify-between items-center gap-x-2">
                <a className="flex-none font-semibold text-xl text-black dark:text-white" href="#" aria-label="Brand">
                  Ekstra
                </a>
                <button
                  onClick={() => setIsMobileSidebarOpen(false)}
                  className="p-1 rounded-md hover:bg-gray-200 dark:hover:bg-neutral-700 hover:cursor-pointer"
                  aria-label="Close sidebar"
                >
                  <XMarkIcon className="h-6 w-6 cursor-pointer" />
                </button>
              </header>
              
              <nav className="flex-1 overflow-y-auto pt-1.5">
                <div className="p-2 w-full flex flex-col flex-wrap">
                  <ul className="space-y-1">
                    {menuItems.map((item) => (
                      <li key={item.name}>
                        <a
                          className={`w-full flex items-center gap-x-3.5 p-2.5 text-base rounded-lg hover:bg-gray-100 dark:hover:bg-neutral-700 ${
                            activeItem === item.name ? "bg-gray-100 dark:bg-neutral-700" : ""
                          }`}
                          href={item.href}
                          onClick={() => handleItemClick(item.name)}
                        >
                          {item.icon}
                          {item.name}
                        </a>
                      </li>
                    ))}
                  </ul>
                </div>
              </nav>
              
              {/* Mobile Account Info Section */}
              <footer className="p-4 border-t border-gray-200 dark:border-neutral-700">
                <AccountDropdown />
              </footer>
            </div>
          </aside>
        </div>
      )}
    </>
  );
}