CREATE TABLE IF NOT EXISTS log (
    Id Int32,
    CampaignId Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    EventTime DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(EventTime)
ORDER BY (EventTime, CampaignId, Id)
SETTINGS index_granularity = 8192;