import React from 'react';

const MainLayout = ({ children, title }) => {
  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <div className="container mx-auto p-4">
        <h1 className="text-3xl font-bold text-center mb-8">{title}</h1>
        {children}
      </div>
    </div>
  );
};

export default MainLayout;