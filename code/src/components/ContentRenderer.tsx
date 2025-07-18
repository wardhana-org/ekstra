import React from 'react';

type ContentItem = {
  type: string;
  url: string;
  title: string;
  size: string;
};

export default function ContentRenderer( {contentItem} : {contentItem : ContentItem[]} ) {

  function getFileIcon(type: string): React.ReactNode {
    switch (type) {
      case 'image':
        return (
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
          </svg>
        );
      case 'video':
        return (
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v8a2 2 0 01-2 2h-2a2 2 0 01-2-2V6z" />
          </svg>
        );
      case 'document':
        return (
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clipRule="evenodd" />
          </svg>
        );
      default:
        return null;
    }
  };

  return (
    <div className="mt-4 space-y-3">
      {contentItem.map((item, index) => (
        <div key={index} className="flex items-center p-3 bg-neutral-700 rounded-lg border border-neutral-600">
          <div className="flex-shrink-0 p-2 bg-neutral-800 rounded-md mr-3">
            {getFileIcon(item.type)}
          </div>
          <div className="flex-grow min-w-0">
            <p className="text-sm font-medium text-white truncate">{item.title}</p>
            <p className="text-xs text-neutral-400">{item.size}</p>
          </div>
          <a 
            href={item.url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="ml-3 px-3 py-1 text-sm bg-neutral-600 hover:bg-neutral-500 rounded-md transition-colors flex items-center"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
            View
          </a>
        </div>
      ))}
    </div>
  );
};