import React from 'react';
import { Table } from 'antd';
import type { TableProps } from 'antd';

interface DataTableProps<T> extends TableProps<T> {
  // Extends Ant Design Table props
}

function DataTable<T extends object>(props: DataTableProps<T>) {
  return (
    <Table
      {...props}
      scroll={{ x: 'max-content' }}
      pagination={{
        showSizeChanger: true,
        showTotal: (total) => `Total ${total} items`,
        ...props.pagination,
      }}
    />
  );
}

export default DataTable;
