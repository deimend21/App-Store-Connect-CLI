# Asset Library items in review submissions

## Contract and evidence

Extend the existing generic review submission commands with `appAssetLibraryImages` and `appAssetLibraryVideos`. No new command group or orchestration wrapper is needed. `asc review items-add` and `asc review items add` previously rejected these resource types.

As observed on October 6, 2026, the published OpenAPI snapshot does not describe these new relationships. The observed App Store Connect browser bundle maps image and video resources to `appAssetLibraryImage` and `appAssetLibraryVideo` review item relationships. Its review item POST uses the existing `/v1/reviewSubmissionItems` endpoint with `reviewSubmission` linkage plus one asset linkage whose resource type is the plural form. An authenticated public GET of `/v1/reviewSubmissions/{id}/items` returned HTTP 200 with `fields[reviewSubmissionItems]=state,appAssetLibraryImage,appAssetLibraryVideo` and `include=appAssetLibraryImage,appAssetLibraryVideo`. This verifies read support and the browser request contract; it does not prove a public POST or App Review acceptance. The official OpenAPI snapshot remains unchanged.

## Workflow

Apple's [Submit assets from Asset Library](https://developer.apple.com/help/app-store-connect/manage-submissions-to-app-review/submit-assets-from-asset-library) help describes standalone submission after an approved app version and co-submission with the first app version. For an app with an approved version, create a review submission and attach the assets independently:

```sh
asc review submissions-create --app APP_ID --platform IOS
asc review items-add --submission SUBMISSION_ID --item-type appAssetLibraryImages --item-id IMAGE_ID
asc review items-add --submission SUBMISSION_ID --item-type appAssetLibraryVideos --item-id VIDEO_ID
asc review items-list --submission SUBMISSION_ID --include appAssetLibraryImage,appAssetLibraryVideo
asc review submissions-submit --id SUBMISSION_ID --confirm
```

For the first app version, submit the assets together with the version. Assets assigned to an active parent submission cannot be submitted separately. Apple validates asset and submission readiness; adding an item alone does not submit it. The existing submission confirmation requirement remains in place. No actual submission was performed for this change.

## Compatibility and validation

The additional types are additive. Existing flags, timeouts, output formats, errors, and destructive-operation confirmations retain their behavior. List responses retain asset relationships and included resource envelopes; table and Markdown show their type and ID. Review history includes asset targets without requiring an app version. `--if-exists skip` reads back the matching relationship and retains strict collection/linkage validation: wrong resource types, null linkages, removed items, and unrelated conflicts cannot turn failure into success.

CLI tests establish RED by rejecting the new types, then check exact image/video POST bodies, list fields and includes, JSON included content, table/Markdown targets, history targets, strict duplicate read-back, and unknown-type validation before authentication. Existing submission tests cover required `--confirm`. Provider submission and App Review acceptance remain unverified.
