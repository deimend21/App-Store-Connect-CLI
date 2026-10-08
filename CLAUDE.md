@AGENTS.md

# App Store Connect CLI - Practical App Submission Quick Reference

> Detailed end-to-end SOP & Checklist: See [`APP_SUBMISSION_CHECKLIST.md`](./APP_SUBMISSION_CHECKLIST.md)

## Critical Pitfalls & Operational Checklist

1. **Privacy Policy & Support URLs (Must Be Public)**:
   - Must be live, public HTTP/HTTPS URLs (e.g. public Notion page, GitHub Pages, custom domain).
   - **Never** use private GitHub repository URLs (returns 404 to Apple reviewers, triggering Guideline 5.1.1 rejection).
   - URLs can be updated via:
     - Version support URL: `asc localizations update --id "<VERSION_LOCALIZATION_ID>" --support-url "<URL>"`
     - App privacy policy URL: `asc localizations update --type app-info --id "<APP_INFO_LOCALIZATION_ID>" --privacy-policy-url "<URL>"`

2. **App Privacy Declarations (Data Usages) - Web Publish Required**:
   - Apple's REST API cannot legally publish developer privacy questionnaires.
   - If `asc review submit` fails with `Associated errors for /v1/appDataUsages/: You must have published answers to your app's data usages`:
     - Visit `https://appstoreconnect.apple.com/apps/<APP_ID>/appPrivacy`
     - Select "No, we do not collect data from this app" (for offline/non-tracking apps) and click **Publish**.

3. **App Review Contact Details**:
   - Phone numbers must include international dialing prefix (e.g. `+86 152xxxxxxxx` instead of `152xxxxxxxx`).
   - Specify `demoAccountRequired: false` if no login is required, or supply valid demo credentials.

4. **Screenshot Dimensions (Pixel Perfect)**:
   - iPhone 6.7": `1290 x 2796` (portrait) or `2796 x 1290` (landscape).
   - iPad 12.9" (if Universal): `2048 x 2732` (portrait) or `2732 x 2048` (landscape).
   - No transparent/alpha channels.

5. **Project Build & Info.plist Requirements**:
   - `CFBundleDisplayName`: User-facing name under the home screen icon.
   - `UILaunchScreen: {}` / Launch Storyboard required for fullscreen compliance.
   - `UIRequiresFullScreen: true` if iPad multitasking is not supported.

6. **Build Lifecycle**:
   - After `asc builds upload`, Apple performs automated processing (5~15 min).
   - Verify processing state is `VALID` (`asc builds list --sort "-uploadedDate" --limit 1`) before submitting.
