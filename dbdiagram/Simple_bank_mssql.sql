CREATE TABLE [accounts] (
  [id] bigserial PRIMARY KEY,
  [owner] nvarchar(255) NOT NULL,
  [blance] bigint NOT NULL,
  [currency] nvarchar(255) NOT NULL,
  [created_at] timestamptz DEFAULT 'now()'
)
GO

CREATE TABLE [entries] (
  [id] bigserial PRIMARY KEY,
  [account_id] bigint NOT NULL,
  [amount] bigint NOT NULL,
  [created_at] timestamptz DEFAULT 'now()'
)
GO

CREATE TABLE [transfers] (
  [id] bigserial PRIMARY KEY,
  [from_account_id] bigint NOT NULL,
  [to_account_id] bigint NOT NULL,
  [amount] bigint NOT NULL,
  [created_at] timestamptz DEFAULT 'now()'
)
GO

ALTER TABLE [entries] ADD FOREIGN KEY ([account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([from_account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([to_account_id]) REFERENCES [accounts] ([id])
GO

CREATE INDEX [accounts_index_0] ON [accounts] ("owner")
GO

CREATE INDEX [entries_index_1] ON [entries] ("account_id")
GO

CREATE INDEX [transfers_index_2] ON [transfers] ("from_account_id")
GO

CREATE INDEX [transfers_index_3] ON [transfers] ("to_account_id")
GO

CREATE INDEX [transfers_index_4] ON [transfers] ("from_account_id", "to_account_id")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '正数金额表示转入，负数表示转出',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entries',
@level2type = N'Column', @level2name = 'amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '转账金额只能是正数',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'transfers',
@level2type = N'Column', @level2name = 'amount';
GO
