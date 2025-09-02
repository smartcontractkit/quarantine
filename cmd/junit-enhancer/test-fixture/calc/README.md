### Notes

This is here to purposefully cause a failure in TestMain. This is because we need to filter out these "test cases".

These only appear when `TestMain` fails, and that can be related to downstream changes.

### Standard JUnit Entry

Here is the standard JUnit entry for a failure in `TestMain` for a package. Notice: `classname="" name="TestMain"classname="" name="TestMain"`.

```xml
	<testsuite tests="1" failures="0" time="0.183000" name="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/calc" timestamp="2025-09-02T15:20:16-07:00">
		<properties>
			<property name="go.version" value="go1.25.0 darwin/arm64"></property>
		</properties>
		<testcase classname="" name="TestMain" time="0.000000">
			<failure message="Failed" type="">[fixture] global setup&#xA;PASS&#xA;[fixture] global teardown&#xA;[fixture] failing on purpose from TestMain&#xA;FAIL&#x9;github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/calc&#x9;0.183s&#xA;</failure>
		</testcase>
		<testcase classname="github.com/smartcontractkit/quarantine/cmd/junit-enhancer/test-fixture/calc" name="TestAdd" time="0.000000"></testcase>
	</testsuite>
```


