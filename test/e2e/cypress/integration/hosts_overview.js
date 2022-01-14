import { allHostNames, agents } from '../fixtures/hosts-overview/available_hosts'

context('Hosts Overview', () => {
    const availableHosts = allHostNames()
    before(() => {
        cy.task('startAgentHeartbeat', agents())
        cy.visit('/');
        cy.navigateToItem('Hosts')
        cy.url().should('include', '/hosts');
    })

    describe('Registered Hosts should be available in the overview', () => {
        it('should show 10 of the 27 registered hosts with default pagination settings', () => {
            cy.get('.tn-hostname').its('length').should('eq', 10)
        })
        it('should show all of the all 27 registered hosts when increasing pagination limit to 100', () => {
            cy.reloadList('hosts', 100)
            cy.get('.tn-hostname').its('length').should('eq', 27)
        })
        describe('Discovered hostnames are the expected ones', () => {
            availableHosts.forEach((hostName) => {
                it(`should have a host named ${hostName}`, () => {
                    cy.get('.tn-hostname a').each(($link) => {
                        const displayedHostName = $link.text().trim()
                        expect(availableHosts).to.include(displayedHostName)
                    })
                })
            })
        })
    })

    describe('Health Detection', () => {
        describe('Health Container shows the health overview of the entire cluster', () => {
            it('should show health status of the first 10 visible hosts', () => {
                cy.log("This test needs to be removed in favor of having the Health Container showing information of the entire cluster")
                cy.reloadList('hosts', 10)
                cy.get('.health-container .health-passing').should('contain', 10)
            })
            it('should show health status of the entire cluster of 27 hosts', () => {
                cy.reloadList('hosts', 100)
                cy.get('.health-container .health-passing').should('contain', 27)
            })
        })
    
        describe('Detected hosts Health matches deployed server status', () => {
            it('all 27 hosts in the cluster should be up', () => {
                availableHosts.forEach((hostName) => cy.get(`#host-${hostName} > .row-status > i`).should('have.class', 'text-success'))
            })
        })
    })

    describe('Hosts Tagging', () => {
        before(() => {
            cy
            .get('body')
            .then(($body) => {
                const deleteTag = '.tn-host-tags x'
                if ($body.find(deleteTag).length > 0) {
                    cy.get(deleteTag).then(($deleteTag) => cy.wrap($deleteTag).click({ multiple: true }))
                }
            })
        })
        const hostsByMatchingPattern = (pattern) => (hostName) => hostName.includes(pattern)
        const taggingRules = [
            ['prd', 'env1'],
            ['qas', 'env2'],
            ['dev', 'env3'],
        ]
        taggingRules.forEach(([pattern, tag]) => {
            describe(`Add tag '${tag}' to all hosts with '${pattern}' in the hostname`, () => {
                availableHosts
                .filter(hostsByMatchingPattern(pattern))
                .forEach((hostName) => {
                    it(`should tag host '${hostName}'`, () => {
                        cy.get(`#host-${hostName} > .tn-host-tags > .tagify`).type(tag).trigger('change')
                    })
                })
            })
        })
    })

    describe('Filtering the Host overview', () => {
        before(() => {
            cy.reloadList('hosts', 100)
        })

        const resetFilter = (option) => {
            cy.intercept('GET', `/hosts?per_page=100`).as('resetFilter')
            cy.get(option).click()
            cy.wait('@resetFilter')
        }

        describe('Filtering by health', () => { 
            before(() => {
                cy.get('.tn-filters > :nth-child(2) > .btn').click()
            })
            const healthScenarios = [
                ['passing', 27],
                ['warning', 0],
                ['critical', 0],
            ]
            healthScenarios.forEach(([health, expectedHostsWithThisHealth], index) => {
                it(`should show ${expectedHostsWithThisHealth || 'an empty list of'} hosts when filtering by health '${health}'`, () => {
                    cy.intercept('GET', `/hosts?per_page=100&health=${health}`).as('filterByHealthStatus') 
                    const selectedOption = `#bs-select-1-${index}`
                    cy.get(selectedOption).click()
                    cy.wait('@filterByHealthStatus').then(() => {
                        expectedHostsWithThisHealth == 0 && cy.get('.table.eos-table').contains('There are currently no records to be shown')
                        expectedHostsWithThisHealth > 0 && cy.get('.tn-hostname').its('length').should('eq', expectedHostsWithThisHealth)
                        resetFilter(selectedOption)
                    })
                })
            })
        })

        describe('Filtering by SAP system', () => {
            before(() => {
                cy.get('.tn-filters > :nth-child(3) > .btn').click()
            })
            const SAPSystemsScenarios = [
                ['HDD', 2],
                ['HDP', 2],
                ['HDQ', 2],
                ['NWD', 4],
                ['NWP', 4],
                ['NWQ', 4],
                
            ]
            SAPSystemsScenarios.forEach(([sapsystem, expectedRelatedHosts], index) => {
                it(`should have ${expectedRelatedHosts} hosts related to SAP system '${sapsystem}'`, () => {
                    cy.intercept('GET', `/hosts?per_page=100&sids=${sapsystem}`).as('filterBySAPSystem')            
                    const selectedOption = `#bs-select-2-${index}`
                    cy.get(selectedOption).click()
                    cy.wait('@filterBySAPSystem').then(() => {
                        cy.get('.tn-hostname').its('length').should('eq', expectedRelatedHosts)
                        resetFilter(selectedOption)
                    })
                })
            })
        })

        describe('Filtering by tags', () => {
            before(() => {
                cy.get('.tn-filters > :nth-child(4) > .btn').click()
            })
            const tagsScenarios = [
                ['env1', 8],
                ['env2', 8],
                ['env3', 8]
            ]
            tagsScenarios.forEach(([tag, expectedTaggedHosts], index) => {
                it(`should have ${expectedTaggedHosts} hosts tagged with tag '${tag}'`, () => {
                    cy.intercept('GET', `/hosts?per_page=100&tags=${tag}`).as('filterByTags')
                    const selectedOption = `#bs-select-3-${index}`
                    cy.get(selectedOption).click()
                    cy.wait('@filterByTags').then(() => {
                        cy.get('.tn-hostname').its('length').should('eq', expectedTaggedHosts)
                        resetFilter(selectedOption)
                    })
                })
            })
        })
    })
});