# Test Fixture

This contains modules, sub-modules, and packages to hit a variety of edge cases.


## math/, utils/strings/

These are part of the main test-fixture sub module. These are happy-case packages which should work, and lead to proper test to file mappings.

<details>

<summary>Example JUnit (before enhance)</summary>

```xml
<testsuite tests="14" failures="0" time="0.366000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" timestamp="2025-10-21T16:04:06-07:00">
  <properties>
    <property name="go.version" value="go1.24.7 darwin/arm64"></property>
  </properties>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestDivide" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestDivideByZero" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestPower/#00" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestPower/#01" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestPower/#02" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestPower/#03" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestPower" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestFactorial/base_cases" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestFactorial/positive_numbers" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="TestFactorial" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="FuzzDivide/seed#0" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="FuzzDivide/seed#1" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="FuzzDivide/seed#2" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/math" name="FuzzDivide" time="0.000000"></testcase>
</testsuite>

<testsuite tests="21" failures="0" time="0.754000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" timestamp="2025-10-21T16:04:06-07:00">
  <properties>
    <property name="go.version" value="go1.24.7 darwin/arm64"></property>
  </properties>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse/hello" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse/#00" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse/a" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse/12345" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse/Hello_World" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestReverse" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestIsPalindrome/palindromes" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestIsPalindrome/not_palindromes" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestIsPalindrome" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#00" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#01" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#02" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#03" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#04" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords/#05" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="TestCountWords" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="FuzzReverse/seed#0" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="FuzzReverse/seed#1" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="FuzzReverse/seed#2" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="FuzzReverse/seed#3" time="0.000000"></testcase>
  <testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/utils/strings" name="FuzzReverse" time="0.000000"></testcase>
</testsuite>
```

</details>

## service/

This is a happy-case sub-module.

<details>

<summary>Example JUnit (before enhance)</summary>


```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="48" failures="0" errors="0" time="0.513166">
	<testsuite tests="14" failures="0" time="0.513000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" timestamp="2025-10-21T16:14:00-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_CreateUser" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_GetUser/non_existent_user" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_GetUser/existing_user" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_GetUser" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_DeleteUser/non_existent_user" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_DeleteUser/existing_user" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_DeleteUser" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_ListUsers/empty_list" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_ListUsers/with_users" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="TestUserService_ListUsers" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="FuzzCreateUser/seed#0" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="FuzzCreateUser/seed#1" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="FuzzCreateUser/seed#2" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service" name="FuzzCreateUser" time="0.000000"></testcase>
	</testsuite>
	<testsuite tests="12" failures="0" time="0.213000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" timestamp="2025-10-21T16:14:00-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_GenerateToken" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_ValidateToken/valid_token" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_ValidateToken/invalid_token" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_ValidateToken/expired_token" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_ValidateToken" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_RevokeToken/existing_token" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_RevokeToken/non_existent_token" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="TestTokenManager_RevokeToken" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="FuzzGenerateToken/seed#0" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="FuzzGenerateToken/seed#1" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="FuzzGenerateToken/seed#2" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/auth" name="FuzzGenerateToken" time="0.000000"></testcase>
	</testsuite>
	<testsuite tests="22" failures="0" time="0.344000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" timestamp="2025-10-21T16:14:00-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Create/valid_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Create/empty_name" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Create/negative_price" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Create" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByID/existing_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByID/non_existent_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByID" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByCategory/electronics_category" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByCategory/furniture_category" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByCategory/non_existent_category" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_GetByCategory" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_UpdateStock/update_existing_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_UpdateStock/update_non_existent_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_UpdateStock" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Delete/delete_existing_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Delete/delete_non_existent_product" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="TestProductRepository_Delete" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="FuzzProductRepository_Create/seed#0" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="FuzzProductRepository_Create/seed#1" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="FuzzProductRepository_Create/seed#2" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="FuzzProductRepository_Create/seed#3" time="0.000000"></testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/service/database/models" name="FuzzProductRepository_Create" time="0.000000"></testcase>
	</testsuite>
</testsuites>
```

</details>

## timeout/

Timeout package explicitly sleeps during a test which causes a timeout error. This will filter out the TestMain timeout,
and only leave the failure for the test that actually timed out. It also contains a test that will work successfully, so
that test should be "enhanced"

**Notes**
* When a test times out, it's possible that it's execution was blocking other tests from starting. So, once the timeout
has been reached and the timeout is reported, we fail assuming that not all tests were able to run.
* I also tried to make a module which had only a `TestMain` which slept longer than the test timeout,
but this didn't register any tests or failures, so it's not necessary.

<details>

<summary>Example JUnit (before enhance)</summary>

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="2" failures="2" errors="0" time="10.232554">
	<testsuite tests="2" failures="1" time="10.233000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/timeout" timestamp="2025-10-22T09:57:04-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="" name="TestMain" time="0.000000">
			<failure message="Failed" type="">panic: test timed out after 10s&#xA;&#x9;running tests:&#xA;&#x9;&#x9;TestCurrentTime_ShouldTimeout (10s)&#xA;&#xA;goroutine 20 [running]:&#xA;testing.(*M).startAlarm.func1()&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:2484 +0x308&#xA;created by time.goFunc&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/time/sleep.go:215 +0x38&#xA;&#xA;goroutine 1 [chan receive]:&#xA;testing.tRunner.func1()&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:1753 +0x43c&#xA;testing.tRunner(0x14000002fc0, 0x1400010cc68)&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:1798 +0x120&#xA;testing.runTests(0x1400000c018, {0x10285b0a0, 0x2, 0x2}, {0x1400007c140?, 0x7?, 0x1028637c0?})&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:2277 +0x3ec&#xA;testing.(*M).Run(0x1400019a0a0)&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:2142 +0x588&#xA;main.main()&#xA;&#x9;_testmain.go:47 +0x90&#xA;&#xA;goroutine 3 [sleep]:&#xA;time.Sleep(0x37e11d600)&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/runtime/time.go:338 +0x158&#xA;github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/timeout.TestCurrentTime_ShouldTimeout(0x14000003180?)&#xA;&#x9;/Users/erik/Documents/repos/quarantine/cmd/junit-enhancer/test-fixture/timeout/timeout_test.go:10 +0x2c&#xA;testing.tRunner(0x14000003180, 0x102776568)&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:1792 +0xe4&#xA;created by testing.(*T).Run in goroutine 1&#xA;&#x9;/Users/erik/.asdf/installs/golang/1.24.7/go/src/testing/testing.go:1851 +0x374&#xA;FAIL&#x9;github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/timeout&#x9;10.231s&#xA;</failure>
		</testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/timeout" name="TestCurrentTime_ShouldTimeout" time="-0.000000">
			<failure message="Failed" type="">=== RUN   TestCurrentTime_ShouldTimeout&#xA;=== PAUSE TestCurrentTime_ShouldTimeout&#xA;=== CONT  TestCurrentTime_ShouldTimeout&#xA;</failure>
		</testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/timeout" name="TestCurrentTime_ShouldNotTimeout" time="0.000000"></testcase>
	</testsuite>
</testsuites>
```

</details>

## buildfailure/

This module contains a build error in the test file which will cause junit enhancer to fail (exit code 1). It also contains
a sub-package that has a valid test, which should be "enhanced" even if a build failure exists for a separate package.

<details>

<summary>Example JUnit (before enhance)</summary>

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="1" failures="1" errors="3" time="0.263308">
	<testsuite tests="0" failures="0" time="0.000000" name="" timestamp="0001-01-01T00:00:00Z">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
	</testsuite>
	<testsuite tests="0" failures="0" time="0.000000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/buildfailure" timestamp="2025-10-22T09:53:22-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="" name="TestMain" time="0.000000">
			<failure message="Failed" type="">FAIL&#x9;github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/buildfailure [build failed]&#xA;</failure>
		</testcase>
	</testsuite>
	<testsuite tests="1" failures="0" time="0.202000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/buildfailure/compiles" timestamp="2025-10-22T09:53:22-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/buildfailure/compiles" name="TestNoBuildFailure" time="0.000000"></testcase>
	</testsuite>
</testsuites>
```

</details>

## testmainfailure/

This contains a `TestMain` definition that purposefully fails.

Because the test suite reports no failures, but TestMain has a failure, we should still note a failure (exit code 1).
It should still also enhance the tests that completed without errors.

<details>

<summary>Example JUnit (before enhance)</summary>

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="1" failures="1" errors="0" time="0.220562">
	<testsuite tests="1" failures="0" time="0.221000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/testmainfailure" timestamp="2025-10-21T16:25:11-07:00">
		<properties>
			<property name="go.version" value="go1.24.7 darwin/arm64"></property>
		</properties>
		<testcase classname="" name="TestMain" time="0.000000">
			<failure message="Failed" type="">[fixture] global setup&#xA;PASS&#xA;[fixture] global teardown&#xA;[fixture] failing on purpose from TestMain&#xA;FAIL&#x9;github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/testmainfailure&#x9;0.219s&#xA;</failure>
		</testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/testmainfailure" name="TestAdd" time="0.000000"></testcase>
	</testsuite>
</testsuites>
```

</details>
